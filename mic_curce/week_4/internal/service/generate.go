package service

//почему то не имплементится моки файл в папку mocks, файл кладется рядом

//go:generate powershell -Command "if (Test-Path mocks) { Remove-Item -Recurse -Force mocks }; New-Item -ItemType Directory -Force mocks"
//go:generate minimock -i NoteService -o ./mocks/note_service_minimock.go
