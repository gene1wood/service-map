// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com

package main

import (
	"encoding/json"
	"testing"
)

func TestGetAsset(t *testing.T) {
	op := opContext{}
	op.newContext(dbconn, false, "127.0.0.1")
	// Tests the first asset in service1
	a, err := getAsset(op, 1)
	if err != nil {
		t.Fatalf("getAsset: %v", err)
	}
	if a.Name != "testhost1.mozilla.com" {
		t.Fatalf("getAsset: unexpected asset name")
	}
	if a.Type != "hostname" {
		t.Fatalf("getAsset: unexpected asset type")
	}
	if a.AssetGroupID != 1 {
		t.Fatalf("getAsset: unexpected asset group id")
	}
	if a.Owner.Operator != "operator" {
		t.Fatalf("getAsset: unexpected asset operator")
	}
	if a.Owner.Team != "testservice" {
		t.Fatalf("getAsset: unexpected asset team")
	}
	if a.Owner.TriageKey != "operator-testservice" {
		t.Fatalf("getAsset: unexpected asset triage key")
	}

	// We should have two indicators here
	if len(a.Indicators) != 2 {
		t.Fatalf("getAsset: unexpected number of indicators")
	}
	expectDetails := "{\"noop\":\"no details in test indicator\"}"
	buf, err := json.Marshal(a.Indicators[0].Details)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if string(buf) != expectDetails {
		t.Fatalf("getAsset: unexpected indicator details")
	}
}

func TestGetAssetHostname(t *testing.T) {
	op := opContext{}
	op.newContext(dbconn, false, "127.0.0.1")
	// Tests the first asset in service1
	alist, err := getAssetHostname(op, "testhost1.mozilla.com")
	if err != nil {
		t.Fatalf("getAssetHostname: %v", err)
	}
	if len(alist) != 1 {
		t.Fatalf("getAssetHostname: unexpected number of assets returned")
	}
	a := alist[0]
	if a.Name != "testhost1.mozilla.com" {
		t.Fatalf("getAsset: unexpected asset name")
	}
	if a.Type != "hostname" {
		t.Fatalf("getAsset: unexpected asset type")
	}
	if a.AssetGroupID != 1 {
		t.Fatalf("getAsset: unexpected asset group id")
	}
	if a.Owner.Operator != "operator" {
		t.Fatalf("getAsset: unexpected asset operator")
	}
	if a.Owner.Team != "testservice" {
		t.Fatalf("getAsset: unexpected asset team")
	}
	if a.Owner.TriageKey != "operator-testservice" {
		t.Fatalf("getAsset: unexpected asset triage key")
	}
}

func TestIndicatorsFromEventSource(t *testing.T) {
	op := opContext{}
	op.newContext(dbconn, false, "127.0.0.1")
	alist, err := indicatorsFromEventSource(op, "testing")
	if err != nil {
		t.Fatalf("indicatorsFromEventSource: %v", err)
	}
	if len(alist) != 5 {
		t.Fatalf("indicatorsFromEventSource: unexpected number of assets returned")
	}
	// For each returned element, we should have a single indicator of the correct
	// source type
	for _, x := range alist {
		if len(x.Indicators) != 1 {
			t.Fatalf("indicatorsFromEventSource: unexpected number of indicators returned")
		}
		if x.Indicators[0].EventSource != "testing" {
			t.Fatalf("indicatorsFromEventSource: indicator had incorrect name")
		}
	}
	// Try marshalling it
	_, err = json.Marshal(&alist)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	alist, err = indicatorsFromEventSource(op, "secondeventsource")
	if err != nil {
		t.Fatalf("indicatorsFromEventSource: %v", err)
	}
	if len(alist) != 1 {
		t.Fatalf("indicatorsFromEventSource: unexpected number of assets returned")
	}
	// For each returned element, we should have a single indicator of the correct
	// source type
	for _, x := range alist {
		if len(x.Indicators) != 1 {
			t.Fatalf("indicatorsFromEventSource: unexpected number of indicators returned")
		}
		if x.Indicators[0].EventSource != "secondeventsource" {
			t.Fatalf("indicatorsFromEventSource: indicator had incorrect name")
		}
	}

	// Try requesting an indicator we should have no data for
	alist, err = indicatorsFromEventSource(op, "nonexistent")
	if err != nil {
		t.Fatalf("indicatorsFromEventSource: %v", err)
	}
	if len(alist) != 0 {
		t.Fatalf("indicatorsFromEventSource: should have had no indicators")
	}
}
