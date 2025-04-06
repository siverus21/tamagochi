@echo off
REM Создаем корневую папку проекта и переходим в неё
mkdir tamagotchi
cd tamagotchi

REM Создаем файлы в корне проекта
type nul > docker-compose.yml
type nul > Dockerfile.backend
type nul > Dockerfile.frontend

REM Создаем структуру для backend
mkdir backend
mkdir backend\cmd
mkdir backend\internal
mkdir backend\internal\config
mkdir backend\internal\handler
mkdir backend\internal\middleware
mkdir backend\internal\model
mkdir backend\internal\repository
mkdir backend\internal\service
mkdir backend\pkg
REM Файлы backend
type nul > backend\cmd\main.go
type nul > backend\go.mod
type nul > backend\go.sum

REM Создаем структуру для frontend
mkdir frontend
mkdir frontend\public
mkdir frontend\src
mkdir frontend\src\components
mkdir frontend\src\pages
mkdir frontend\src\services
REM Файлы frontend
type nul > frontend\package.json
type nul > frontend\webpack.config.js

REM Создаем структуру для базы данных
mkdir db
type nul > db\init.sql
mkdir db\data

echo Структура проекта успешно создана!
pause
