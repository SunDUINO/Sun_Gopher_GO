# Gopher Go! 🎮

Platformówka napisana w Go z użyciem Ebitengine.

## Wymagania

- **Go 1.22+** – https://go.dev/dl/
- **Ebitengine** – zostanie pobrane automatycznie przez `go mod tidy`
- Zależności systemowe (Linux): `libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev libasound2-dev`

## Instalacja i uruchomienie

```bash
# 1. Wejdź do katalogu projektu
cd giana-go

# 2. Pobierz zależności
go mod tidy

# 3. Uruchom grę
go run .
```

### Na Linux (jeśli brakuje bibliotek systemowych)
```bash
sudo apt-get install -y libx11-dev libxcursor-dev libxrandr-dev \
  libxinerama-dev libxi-dev libgl1-mesa-dev libasound2-dev
go mod tidy
go run .
```

### Na macOS
```bash
go mod tidy
go run .
```

### Na Windows
```bash
go mod tidy
go run .
```

### UWAGA !!  Może być konieczna edycja go.mod i importu w main jeśli sklonujesz do innego projektu


## Sterowanie

| Klawisz | Akcja |
|---------|-------|
| ← → / A D | Ruch lewo/prawo |
| ↑ / W / Spacja | Skok |
| ESC | Menu główne |
| Enter | Start / Restart |

## Mechaniki gry

- 🪙 **Monety** – zbieraj dla punktów (10 pkt każda)
- 🍄 **Wrogowie** – wejdź na głowę wroga (skocz na niego), żeby go pokonać
- ⚠️ **Kolce** – natychmiastowa śmierć!
- 🏆 **Cel (G)** – żółty portal – dotknij go aby wygrać poziom
- ❤️ **Życia** – masz 3 życia, po stracie wszystkich – Game Over

## Struktura projektu

```
sun_gopher_go/
├── go.mod         – moduł Go
├── src/
    ├── main.go 
├── game
    ├── game.go        – główna pętla gry, rysowanie
    ├── player.go      – gracz: ruch, kolizje, animacja
    ├── entities.go    – monety i wrogowie
    ├── level.go       – dane poziomu (mapa kafelków)
    ├── constants.go   – stałe i typy
```

## Rozszerzenia (pomysły)

- Muzyka/dźwięki – `ebiten/audio`
- Więcej poziomów – dodaj kolejne tablice w `level.go`
- Sprite'y z plików PNG zamiast rysowania prymitywami
- Animowane tła (paralaksa)
- Tryb dwóch graczy
EOF
