// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetCommitBranches(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := openRepositoryWithDefaultContext(bareRepo1Path)
	require.NoError(t, err)
	defer bareRepo1.Close()

	// these test case are specific to the repo1_bare test repo
	testCases := []struct {
		CommitID         string
		ExpectedBranches []string
	}{
		{"2839944139e0de9737a044f78b0e4b40d989a9e3", []string{"branch1"}},
		{"5c80b0245c1c6f8343fa418ec374b13b5d4ee658", []string{"branch2"}},
		{"37991dec2c8e592043f47155ce4808d4580f9123", []string{"master"}},
		{"95bb4d39648ee7e325106df01a621c530863a653", []string{"branch1", "branch2"}},
		{"8d92fc957a4d7cfd98bc375f0b7bb189a0d6c9f2", []string{"branch2", "master"}},
		{"master", []string{"master"}},
	}
	for _, testCase := range testCases {
		commit, err := bareRepo1.GetCommit(testCase.CommitID)
		require.NoError(t, err)
		branches, err := bareRepo1.getBranches(commit, 2)
		require.NoError(t, err)
		assert.Equal(t, testCase.ExpectedBranches, branches)
	}
}

func TestGetTagCommitWithSignature(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := openRepositoryWithDefaultContext(bareRepo1Path)
	require.NoError(t, err)
	defer bareRepo1.Close()

	// both the tag and the commit are signed here, this validates only the commit signature
	commit, err := bareRepo1.GetCommit("28b55526e7100924d864dd89e35c1ea62e7a5a32")
	require.NoError(t, err)
	assert.NotNil(t, commit)
	assert.NotNil(t, commit.Signature)
	// test that signature is not in message
	assert.Equal(t, "signed-commit\n", commit.CommitMessage)
}

func TestGetCommitWithBadCommitID(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := openRepositoryWithDefaultContext(bareRepo1Path)
	require.NoError(t, err)
	defer bareRepo1.Close()

	commit, err := bareRepo1.GetCommit("bad_branch")
	assert.Nil(t, commit)
	require.Error(t, err)
	assert.True(t, IsErrNotExist(err))
}

func TestIsCommitInBranch(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := openRepositoryWithDefaultContext(bareRepo1Path)
	require.NoError(t, err)
	defer bareRepo1.Close()

	result, err := bareRepo1.IsCommitInBranch("2839944139e0de9737a044f78b0e4b40d989a9e3", "branch1")
	require.NoError(t, err)
	assert.True(t, result)

	result, err = bareRepo1.IsCommitInBranch("2839944139e0de9737a044f78b0e4b40d989a9e3", "branch2")
	require.NoError(t, err)
	assert.False(t, result)
}

func TestRepository_CommitsBetweenIDs(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo4_commitsbetween")
	bareRepo1, err := openRepositoryWithDefaultContext(bareRepo1Path)
	require.NoError(t, err)
	defer bareRepo1.Close()

	cases := []struct {
		OldID           string
		NewID           string
		ExpectedCommits int
	}{
		{"fdc1b615bdcff0f0658b216df0c9209e5ecb7c78", "78a445db1eac62fe15e624e1137965969addf344", 1}, // com1 -> com2
		{"78a445db1eac62fe15e624e1137965969addf344", "fdc1b615bdcff0f0658b216df0c9209e5ecb7c78", 0}, // reset HEAD~, com2 -> com1
		{"78a445db1eac62fe15e624e1137965969addf344", "a78e5638b66ccfe7e1b4689d3d5684e42c97d7ca", 1}, // com2 -> com2_new
	}
	for i, c := range cases {
		commits, err := bareRepo1.CommitsBetweenIDs(c.NewID, c.OldID)
		require.NoError(t, err)
		assert.Len(t, commits, c.ExpectedCommits, "case %d", i)
	}
}
