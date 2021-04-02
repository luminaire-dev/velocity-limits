package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	main()
	expected, err := os.Open("./output.txt")
	require.NoError(t, err)
	defer expected.Close()

	actual, err := os.Open("./generated_output.txt")
	require.NoError(t, err)
	defer actual.Close()

	scanExp := bufio.NewScanner(expected)
	scanAct := bufio.NewScanner(actual)

	for scanExp.Scan() && scanAct.Scan() {
		require.Equal(t, scanExp.Text(), scanAct.Text())
	}

	err = scanAct.Err()
	require.NoError(t, err)

	err = scanExp.Err()
	require.NoError(t, err)
}
