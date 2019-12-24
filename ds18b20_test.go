package ds18b20

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSensors = []string{"28-000000000000", "28-000000000001", "28-000000000002", "28-000000000003", "28-000000000004", "28-000000000005", "28-000000000006"}

func TestMain(m *testing.M) {
	busPath = "testdata/"
	busMasterPath = busPath + "w1_bus_master1/"
	busMaster = busMasterPath + "w1_master_slaves"

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestGetSensors(t *testing.T) {
	t.Log("Running TestGetSensors")

	testCases := []struct {
		name            string
		busMaster       string
		expectedSensors []*Sensor
		expectedError   string
	}{
		{
			name: "Correct sensors",
			expectedSensors: []*Sensor{
				&Sensor{Name: testSensors[1]},
				&Sensor{Name: testSensors[2]},
				&Sensor{Name: testSensors[3]},
			},
		},
		{
			name:          "wrong bus master",
			busMaster:     "wrongfile",
			expectedError: "open wrongfile: no such file or directory",
		},
	}

	for _, testCase := range testCases {
		if testCase.busMaster != "" {
			busMaster = testCase.busMaster
		}
		sensors, err := GetSensors()

		// assert error
		if testCase.expectedError != "" {
			assert.EqualError(t, err, testCase.expectedError)
			continue
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, testCase.expectedSensors, sensors)
	}
}

func TestGetTemperature(t *testing.T) {
	t.Log("Running TestGetTemperature")
	testCases := []struct {
		name          string
		sensor        *Sensor
		expectedTemp  float64
		expectedError string
	}{
		{
			name:          "Not existing sensor",
			sensor:        &Sensor{Name: testSensors[0]},
			expectedError: `failed to read sensor "28-000000000000": open testdata/28-000000000000/w1_slave: no such file or directory`,
		},
		{
			name:         "Sensor with correct CRC",
			sensor:       &Sensor{Name: testSensors[1]},
			expectedTemp: 20.562,
		},
		{
			name:          "Sensor with incorrect CRC",
			sensor:        &Sensor{Name: testSensors[2]},
			expectedError: "could not verify crc, regex did not match",
		},
		{
			name:          "Sensor with incorrect temperature format",
			sensor:        &Sensor{Name: testSensors[3]},
			expectedError: "could not extract temperature, regex did not match",
		},
		{
			name:          "Sensor with no content",
			sensor:        &Sensor{Name: testSensors[4]},
			expectedError: "can not parse temperature, wrong file format",
		},
		{
			name:          "Sensor with strange temperature value",
			sensor:        &Sensor{Name: testSensors[5]},
			expectedError: "could not extract temperature, parse to float failed",
		},
		{
			name:         "Sensor with below 0 temperature value",
			sensor:       &Sensor{Name: testSensors[6]},
			expectedTemp: -5.32,
		},
	}

	for _, testCase := range testCases {
		t.Logf("- test %q", testCase.name)
		got, err := testCase.sensor.GetTemperature()
		t.Logf("  - error got: %v, error expected: %v\n", err, testCase.expectedError)

		// assert error
		if testCase.expectedError != "" {
			assert.EqualError(t, err, testCase.expectedError)
			continue
		} else {
			assert.NoError(t, err)
		}

		// assert value is correct
		assert.Equal(t, testCase.expectedTemp, got, "Received the expected temperature")
	}
}
