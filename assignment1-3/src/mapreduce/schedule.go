package mapreduce

import (
	"sync"
)

func (mr *Master) dispatchWorker(task int, nOtherPhase int, phase jobPhase, wg *sync.WaitGroup) {
	debug("Starting worker: %v\n", task)
	worker := <-mr.registerChannel

	var args DoTaskArgs
	args.JobName = mr.jobName
	args.File = mr.files[task]
	args.Phase = phase
	args.TaskNumber = task
	args.NumOtherPhase = nOtherPhase

	ok := call(worker, "Worker.DoTask", args, new(struct{}))
	if ok {
		wg.Done()
		mr.registerChannel <- worker
	} else {
		// if the master's RPC to the worker fails,
		// the master should re-assign the task given to the failed worker to another worker.
		debug("Worker failed with task: %v\n", task)
		mr.dispatchWorker(task, nOtherPhase, phase, wg)
	}
}

// _____________________________________________________________________________
// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	var workerGroup sync.WaitGroup
	for task := 0; task < ntasks; task++ {
		workerGroup.Add(1)
		go mr.dispatchWorker(task, nios, phase, &workerGroup)
	}
	workerGroup.Wait()

	//schedule only needs to tell the workers the name of the original
	//input file (mr.files[task]) and the task task; each worker knows
	//from which files to read its input and to which files to write its output.
	//The master tells the worker about a new task by sending it the RPC call Worker.DoTask,
	//giving a DoTaskArgs object as the RPC argument.

	debug("Schedule: %v phase done\n", phase)
}
