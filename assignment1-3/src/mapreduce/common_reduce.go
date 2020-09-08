package mapreduce

import (
	"encoding/json"
	"io"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Use checkError to handle errors.

	keyValsMap := make(map[string][]string)
	for m := 0; m < nMap; m++ {
		fileName := reduceName(jobName, m, reduceTaskNumber)
		interFile, err := os.Open(fileName)
		checkError(err)
		defer interFile.Close()

		dec := json.NewDecoder(interFile)
		for {
			var keyVal KeyValue
			if err := dec.Decode(&keyVal); err == io.EOF {
				break
			} else {
				checkError(err)
			}
			keyValsMap[keyVal.Key] = append(keyValsMap[keyVal.Key], keyVal.Value)
		}
	}

	writeFileName := mergeName(jobName, reduceTaskNumber)
	writeFile, err := os.Create(writeFileName)
	checkError(err)
	defer writeFile.Close()

	var sortedKeys []string
	for k := range keyValsMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	enc := json.NewEncoder(writeFile)
	for _, key := range sortedKeys {
		enc.Encode(KeyValue{key, reduceF(key, keyValsMap[key])})
	}
}
