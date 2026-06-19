package repository

//для линукс sh -c "rm -rf mocks && mkdir -p mocks"

//go:generate powershell -Command "if (Test-Path mocks) { Remove-Item -Recurse -Force mocks }; New-Item -ItemType Directory -Force mocks"
//go:generate minimock -i NoteRepository -o ./mocks/note_repository_minimock.go
