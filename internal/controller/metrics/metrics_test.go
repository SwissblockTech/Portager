/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSyncTotalRegistered(t *testing.T) {
	SyncTotal.WithLabelValues("test-sync", "default", "success").Inc()

	expected := `
		# HELP portage_sync_total Total number of ImageSync reconcile completions
		# TYPE portage_sync_total counter
		portage_sync_total{name="test-sync",namespace="default",status="success"} 1
	`
	if err := testutil.CollectAndCompare(SyncTotal, strings.NewReader(expected)); err != nil {
		t.Errorf("SyncTotal metric mismatch: %v", err)
	}
}

func TestSyncDurationRegistered(t *testing.T) {
	SyncDuration.WithLabelValues("test-sync", "default").Observe(5.0)

	count := testutil.CollectAndCount(SyncDuration)
	if count == 0 {
		t.Error("SyncDuration metric not registered")
	}
}

func TestImagesCopiedRegistered(t *testing.T) {
	ImagesCopied.WithLabelValues("test-sync", "default").Inc()
	ImagesCopied.WithLabelValues("test-sync", "default").Inc()

	val := testutil.ToFloat64(ImagesCopied.WithLabelValues("test-sync", "default"))
	if val != 2 {
		t.Errorf("expected ImagesCopied=2, got %v", val)
	}
}

func TestImagesSkippedRegistered(t *testing.T) {
	ImagesSkipped.WithLabelValues("test-sync", "default").Inc()

	val := testutil.ToFloat64(ImagesSkipped.WithLabelValues("test-sync", "default"))
	if val != 1 {
		t.Errorf("expected ImagesSkipped=1, got %v", val)
	}
}

func TestImagesFailedRegistered(t *testing.T) {
	ImagesFailed.WithLabelValues("test-sync", "default").Inc()

	val := testutil.ToFloat64(ImagesFailed.WithLabelValues("test-sync", "default"))
	if val != 1 {
		t.Errorf("expected ImagesFailed=1, got %v", val)
	}
}

func TestImageInfoRegistered(t *testing.T) {
	ImageInfo.WithLabelValues("test-sync", "default", "3", "1", "4").Set(1)

	val := testutil.ToFloat64(ImageInfo.WithLabelValues("test-sync", "default", "3", "1", "4"))
	if val != 1 {
		t.Errorf("expected ImageInfo=1, got %v", val)
	}
}

func TestSyncTotalLabelCardinality(t *testing.T) {
	// Verify both status values work without panic.
	SyncTotal.WithLabelValues("cardinality-test", "ns", "success").Inc()
	SyncTotal.WithLabelValues("cardinality-test", "ns", "failure").Inc()
}
