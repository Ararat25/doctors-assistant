# 1. Цель проекта

Цель проекта - разработать удобную систему для взаимодействия с искусственным интеллектом, в которой пользователь (медицинский работник) сможет узнать предполагаемый диагноз у пациента с помощью нейронной сети Sber GigaChat (далее Система). Также в системе должна быть реализована возможность узнавать медицинские новости, смотреть теоретический материал и оставлять сообщения на форуме.

# 2. Описание системы 

Система состоит из следующих основных функциональных блоков:

1. Регистрация, аутентификация и авторизация.
2. ДОПОЛНИТЬ Функционал работы с нейронной сетью. ДОПОЛНИТЬ
3. Функционал просмотра новостей.
4. Функционал общения на форуме.

## 2.1 Типы пользователей

Система предусматривает один тип пользователей системы - доктор. Доктор может запрашивать у нейронной сети диагноз, общаться на форуме и смотреть новости.

## 2.2 Регистрация

1. Доктор переходит на страницу регистрации
2. Вводит ФИО, email, пароль и другие данные
3. Нажимает на кнопку регистрации
4. Сервер получает данные и сверяет их с уже существующими
5. Если введенный логин свободен, то сервер заносит данные в БД и отправляет ответ с подтверждением регистрации
6. Если пользователь с таким логином уже существует, то сервер отправляет отрицательный ответ, с сообщением об ошибке

## 2.3 Аутентификация доктора

1. Доктор переходит на страницу аутентификации
2. Вводит логин и пароль
3. Нажимает кнопку войти
4. Если логин и пароль валидны, то отправляется запрос на сервер
5. В противном случае, пользователю сообщается об ошибке и предлагается ввести данные ещё раз
6. Когда на сервер приходят данные с логином и паролем, он сначала находит пользователя по логину, если это не удалось, то отправляется ответ о том, что такого пользователя не существует
7. Затем сервер сравнивает хэши введенного пароля и хранящегося в БД, если они совпадают, то пользователя пропускает на страницу с его аккаунтом, в противном случае отправляется соответствующее сообщение

## 2.4 Функционал работы с нейронной сетью

Система смотрит специальность врача и исходя из этого, сужает набор возможных симптомов и заболеваний. 

Доктор может либо из списка выбирать симптом, либо писать симптом сам. 

Нейронная сеть может уточнять симптомы и в конечном итоге сообщает доктору возможные заболевания и их вероятность. 

Все диалоги с сетью сохраняются и доступны в любой момент времени.

Если искусственному интеллекту требуются результаты каких-либо анализов для определения диагноза, диалог сохраняется с указанием данных пациента и когда придут результаты, врач сможет открыть нужный диалог и ввести новые данные

## 2.5 Функционал общения на форуме

Предоставляет возможность для обмена информацией, обсуждения тем, задавания вопросов и нахождения ответов. 

Пользователи могут создавать темы, комментировать посты, ставить лайки, делиться своим мнением, опытом и знаниями. 

Имеется возможность фильтрации контента по различным критериям.

Возможен поиск по ключевым словам, авторам или категориям.

Настройка предпочтений происходит исходя из специальности врача.

Имеются инструменты для модерации и контроля качества контента, обеспечивая безопасную и плодотворную обстановку для всех участников форума.

## 2.6 Функционал просмотра новостей

Обеспечивает удобный и быстрый доступ к новостным статьям. 

Настройка предпочтений происходит исходя из специальности врача и наиболее часто просматриваемого контента.

Пользователи могут использовать фильтрацию для выбора интересующих категорий новостей.

Имеется возможность настройки уведомлений. 

Пользователи могут комментировать новостные посты, ставить лайки, делиться своим мнением, опытом и знаниями. 

# 3. Предлагаемый стек технологий

- Бэкенд:
  - Go
- Фронтенд:
  - HTML
  - JavaScript
