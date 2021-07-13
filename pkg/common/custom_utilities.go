// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package common

import (
       	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
)

var	TestDataDirectory = "testdata"
var	DefaultTimestamp = "0001-01-01T00:00:00Z"
var	ReplaceTimestampRegExp = "s/\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z\"/\"" + DefaultTimestamp + "\"/"

// IsModifyingStatus returns true if a
// resource's latest status matches one
// of the provided modifying statuses.
func IsModifyingStatus(
	latestStatus *string,
	modifyingStatuses *[]string,
) bool {
	for _, status := range *modifyingStatuses {
		if *latestStatus == status {
			return true
		}
	}
	return false
}

// Checks to see if the contents of a given yaml file, with name stored
// in expectation, matches the given actualYamlByteArray.
func IsYamlEqual(expectation *string, actualYamlByteArray *[]byte) bool {
     	// Get the file name of the expected yaml.
	expectedYamlFileName := TestDataDirectory + "/" + *expectation

	// Build a tmp file for the actual yaml.
	actualYamlFileName := buildTmpFile("actualYaml", *actualYamlByteArray)
	defer os.Remove(actualYamlFileName)
	if "" == actualYamlFileName {
	  fmt.Printf("Could not create temporary actual file.\n")
	  return false
	}

	// Replace Timestamps that would show up as different.
	_, err := exec.Command("sed", "-r", "-i", ReplaceTimestampRegExp, actualYamlFileName).Output()
	if !checkExecCommandError(err) {
	    return false
	}
	
	output,err := exec.Command("diff", "-c", expectedYamlFileName, actualYamlFileName).Output()
	if !checkExecCommandError(err) {
	   return false
	}

	if len(output) > 0 {
	   actualOutput,_ := exec.Command("cat", actualYamlFileName).Output()
	   fmt.Printf("\nActual Yaml File Instead of: " + expectedYamlFileName + "\n" + string(actualOutput) + "\n")
	   fmt.Printf("Diff From Expected:\n" + string(output) + "\n")
	   return false
	}
	return true
}

func buildTmpFile(fileNameBase string, contents []byte) string {
     newTmpFile, err := ioutil.TempFile(TestDataDirectory, fileNameBase)
     if err != nil {
        fmt.Println(err)
	return ""
     }
     if _, err := newTmpFile.Write(contents); err != nil {
        fmt.Println(err)
	return ""
     }
     if err := newTmpFile.Close(); err != nil {
        fmt.Println(err)
	return ""
     }		
     return newTmpFile.Name()
}

func checkExecCommandError(err error) bool {
     if err == nil {
     	return true
     }
     switch err.(type) {
	case *exec.ExitError:
	       // ExitError is expected.
	       return true
	default:
	       // Couldn't run diff.
	       fmt.Printf("Exec Command Error: ")
	       fmt.Println(err)
	       return false
     }
}
