/*
 * ╔════════════════════════════════════════════════════════════════╗
 * ║ sun_gopher_go                                                  ║
 * ║ Plik / File: main_test.go                                      ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ Autor / Author:                                                ║
 * ║   SunRiver                                                     ║
 * ║   Lothar TeaM                                                  ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ GitHub  : sun_gopher_go                                        ║
 * ║ WWW     : https://lothar-team.pl                               ║
 * ║ Forum   : https://forum.lothar-team.pl                         ║
 * ║                                                                ║
 * ║ Licencja / License: MIT                                        ║
 * ║ Rok / Year: 2026                                               ║
 * ╚════════════════════════════════════════════════════════════════╝
 */
package main

import "testing"

func TestMainApp(t *testing.T) {
	result := "SunGo"
	if result != "SunGo" {
		t.Errorf("Oczekiwano SunGo, otrzymano %s", result)
	}
}