# GB_Backend_Level1

Репозиторий для домашних заданий на курсе "Backend-разработка на Go. Уровень 1"

## ДЗ Урок 2

1. Добавить в приложение рассылки даты/времени возможность отправлять клиентам произвольные сообщения из консоли сервера
2. Добавить в приложение чата возможность устанавливать клиентам свой никнейм при подключении к серверу
3. Дополнительное задание
   Реализовать игру “Математика на скорость”: сервер генерирует случайное выражение с двумя операндами, сохраняет ответ, а затем отправляет выражение
   всем клиентам. Первый клиент, отправивший правильный ответ - побеждает, затем генерируется следующее выражение и так далее.
4. Опишите текущие наработки по курсовому проекту.

Критерии приёма задания:
1. Для заданий 1-3 в комментарий приложите ссылки на пулл-реквест(-ы).
2. Для задания 4 опишите текущие наработки. Можно приложить ссылку на репозиторий, коммиты и/или пулл-реквесты.

## ДЗ Урок 3

1. Добавить параметры для фильтрации товаров по диапазону цены
2. Добавить в спецификацию объект Order (заказ), подумать, какие поля у нее должны быть и какие эндпоинты потребуется описать

## ДЗ Урок 4

1. Добавить в пример с файловым сервером из методички возможность получить список всех файлов на сервере (имя, расширение, размер в байтах).
2. С помощью query-параметра, реализовать фильтрацию выводимого списка по расширению (то есть, выводить только .png файлы, или только .jpeg).
3. *Текущая реализация сервера не позволяет хранить несколько файлов с одинаковым названием (т.к. они будут храниться в одной директории на диске). Подумайте, как можно обойти это ограничение?
4. К коду, написанному в рамках заданий 1-3, добавьте тесты с использованием библиотеки httptest. 
