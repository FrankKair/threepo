package main

import (
    "encoding/json"
    "testing"
)

func TestBuildObject(t *testing.T) {
    rows := [][]string{
        {"key", "en", "pt", "sv"},
        {"hello", "hello", "olá", "hej"},
        {"computer", "computer", "computador", "dator"},
    }

    got := buildObject(rows)
    want := map[string]any{
        "hello":    map[string]string{"en": "hello", "pt": "olá", "sv": "hej"},
        "computer": map[string]string{"en": "computer", "pt": "computador", "sv": "dator"},
    }

    assertJSON(t, got, want)
}

func TestBuildObjectNestedKeys(t *testing.T) {
    rows := [][]string{
        {"key", "en", "pt"},
        {"home.title", "Home", "Início"},
        {"home.subtitle", "Welcome", "Bem-vindo"},
        {"auth.login", "Log in", "Entrar"},
    }

    got := buildObject(rows)

    b, _ := json.Marshal(got)
    var parsed map[string]any
    json.Unmarshal(b, &parsed)

    home := parsed["home"].(map[string]any)
    if home["title"].(map[string]any)["en"] != "Home" {
        t.Errorf("home.title.en = %v, want Home", home["title"].(map[string]any)["en"])
    }
    if home["subtitle"].(map[string]any)["pt"] != "Bem-vindo" {
        t.Errorf("home.subtitle.pt = %v, want Bem-vindo", home["subtitle"].(map[string]any)["pt"])
    }

    auth := parsed["auth"].(map[string]any)
    if auth["login"].(map[string]any)["en"] != "Log in" {
        t.Errorf("auth.login.en = %v, want Log in", home["login"].(map[string]any)["en"])
    }
}

func TestBuildObjectEmpty(t *testing.T) {
    got := buildObject([][]string{{"key", "en"}})
    if len(got) != 0 {
        t.Errorf("expected empty map, got %v", got)
    }
}

func assertJSON(t *testing.T, got, want any) {
    t.Helper()
    g, _ := json.Marshal(got)
    w, _ := json.Marshal(want)
    if string(g) != string(w) {
        t.Errorf("\ngot: %s\nwant: %s", g, w)
    }
}
