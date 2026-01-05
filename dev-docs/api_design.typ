= Проектирование API
= Сервис поиска потерянных вещей
=== Основные ресурсы
+ `auth` --- аутентификация
    + `/api/v1/auth/refresh-token` POST --- обновление access-токена
        - `{refresh_token}`
    + `/api/v1/auth/login` POST --- вход в систему
        - `{email: string, password: string}`
    + `/api/v1/auth/register` POST --- регистрация по инвайт-коду
        - `{code: UUID, email: string, password: string, first_name: string, last_name: string, middle_name?: string}`
    + `/api/v1/auth/logout` POST --- выход из аккаунта
    + `/api/v1/auth/reset-password` POST --- запрос на сброс пароля
        - `{email: string}`
    + `/api/v1/auth/reset-password/confirm` POST --- подтверждение сброса
        - `{token: string, new_password: string}`
+ `users` --- пользователи всех ролей
    + `/api/v1/users?role={role_id}&search={search_text}&limit={limit}&offset={offset}` GET --- получить список всех пользователей
    + `/api/v1/users/{user_id}` GET --- получить информацию о пользователе
    + `/api/v1/users/{user_id}` DELETE --- удалить пользователя
    + `/api/v1/users/check-email` POST --- проверка, свободен ли email
        - `{email: string}`
+ `roles` --- роли пользователей
    + `/api/v1/roles` GET --- получить список всех ролей
    + `/api/v1/roles` POST --- добавить новую роль
        - `{name: string}`
    + `/api/v1/roles/{role_id}` DELETE --- удалить существующую роль
    + `/api/v1/roles/{role_id}` PATCH --- изменить название существующей роли
        - `{name: string}`
+ `permissions` --- права доступа
    + `/api/v1/permissions` GET --- получить список всех прав доступа
    + `/api/v1/permissions` POST --- добавить новое право доступа
        - `{name: string}`
    + `/api/v1/permissions/{permission_id}` DELETE --- удалить существующее право доступа
    + `/api/v1/permissions/{permission_id}` PATCH --- изменить название права доступа
        - `{name: string}`
+ `role_permission` --- связь ролей и прав (связующая таблица; многие ко многим)
    + `/api/v1/roles/{role_id}/permissions` GET --- получить список всех прав роли
    + `/api/v1/roles/{role_id}/permissions` POST --- добавить право роли
        - `{permission_id: int}`
    + `/api/v1/roles/{role_id}/permissions/{permission_id}` DELETE --- удалить право у роли
+ `user_role` --- связь пользователей и их ролей (связующая таблица; многие ко многим)
    + `/api/v1/users/{user_id}/roles` GET --- получить список всех ролей пользователя
    + `/api/v1/users/{user_id}/roles` POST --- добавить пользователю роль
        - `{role_id: int}`
    + `/api/v1/users/{user_id}/roles/{role_id}` DELETE --- удалить роль у пользователя
+ `posts` --- объявления о потерянных вещах
    + `/api/v1/posts?author={author: UUID}&verified=&thing_returned_to_owner=&search=&limit=&offset=` GET --- список объявлений
    + `/api/v1/posts/{post_id: UUID}` GET --- детали объявления
    + `/api/v1/posts` POST --- создать объявление
        - `{name: string, description?: string, photo_url?: string}`
    + `/api/v1/posts/{post_id}` PATCH --- обновить своё объявление
        - `{name?: string, description?: string, photo_url?: string}`
    + `/api/v1/posts/{post_id}` DELETE --- удалить своё объявление
    + `/api/v1/posts/{post_id}/verify` PATCH --- верификация объявления (админы сервиса)
        - `{verified: boolean}`
    + `/api/v1/posts/{post_id}/returned-to-owner` PATCH --- моя вещь найдена, `returned_to_owner` устанавливается в `true`
+ `invite_codes` --- пригласительные коды
    + `/api/v1/invite_codes?role_id=` GET --- список актуальных инвайт-кодов
    + `/api/v1/invite_codes/` POST --- создать инвайт-код (суперадмин, админы сервиса)
        - `{role_id}`
    + `/api/v1/invite_codes/{code}` DELETE --- отозвать инвайт-код (суперадмин, админы сервиса)
+ `students` --- данные учеников
    + `/api/v1/students?group_id=&search=text` GET --- список учеников
    + `/api/v1/students/{student_id}` GET --- информация об ученике
    + `/api/v1/students/{student_id}` PATCH --- обновить данные ученика
        - `{group_id?: int}`
    + `/api/v1/students/{student_id}/parents` GET --- получить родителей ученика
+ `parents` --- данные родителей
    + `/api/v1/parents/{parent_id}/students` GET --- дети родителя
    + `/api/v1/parents/{parent_id}/students/{student_id}` POST --- привязать ребёнка
    + `/api/v1/parents/{parent_id}/students/{student_id}` DELETE --- отвязать ребёнка
+ `teachers` --- данные преподавателей
    + `/api/v1/teachers?subject_id=&search=` GET --- список преподавателей
    + `/api/v1/teachers/{teacher_id}` GET --- информация о преподавателе
    + `/api/v1/teachers/{teacher_id}/subjects` GET --- предметы, которые ведёт преподаватель
    + `/api/v1/teachers/{teacher_id}/subjects` POST --- добавить предмет учителю
        - `{subject_id}`
    + `/api/v1/teachers/{teacher_id}/subjects/{subject_id}` DELETE --- удалить предмет
    + `/api/v1/teachers/{teacher_id}` PATCH --- обновить данные учителя
        - `{classroom_id?: int}`
+ `student_groups` --- классы/группы учеников/студентов
    + `/api/v1/student_groups?group_advisor=UUID` GET --- получить список всех групп/классов
    + `/api/v1/student_groups` POST --- создать новую группу
        - `{name, group_advisor}`
    + `/api/v1/student_groups/{group_id}/students` GET --- получить список учеников в группе
    + `/api/v1/student_groups/{group_id}/group_advisor` GET --- получить (классного) руководителя группы/класса
    + `/api/v1/student_groups/{group_id}` PATCH --- обновить данные группы
        - `{name?, group_advisor?}`
    + `/api/v1/student_groups/{group_id}` DELETE --- удалить группу
+ `staff` --- данные сотрудников ОУ
    + `/api/v1/staff?position_id=` GET --- список всех сотрудников
    + `/api/v1/staff/positions` GET --- список всех должностей
    + `/api/v1/staff/{staff_id}` PATCH --- обновить сотрудника
        - `{position_id: int}`
+ `institution_administrators` --- данные представителей администрации ОУ
    + `/api/v1/institution_administrators?position_id=` GET --- список представителей администрации ОУ
    + `/api/v1/institution_administrators/positions` GET --- список должностей
    + `/api/v1/institution_administrators/{institution_administrators_id}` PATCH --- обновить сотрудника администрации ОУ
        - `{position_id: int}`
+ `rooms` --- аудитории, учебные кабинеты
    + `/api/v1/rooms` GET --- список комнат
    + `/api/v1/rooms` POST --- добавить комнату
        - `{name: string}`
    + `/api/v1/rooms/{room_id}` PATCH --- обновить комнату
        - `{name: string}`
    + `/api/v1/rooms/{room_id}` DELETE --- удалить комнату
+ `subjects` --- предметы
    + `/api/v1/subjects` GET --- список всех предметов
    + `/api/v1/subjects` POST --- создать предмет
        - `{name: string}`
    + `/api/v1/subjects/{subject_id}` PATCH --- обновить предмет
        - `{name: string}`
    + `/api/v1/subjects/{subject_id}` DELETE --- удалить предмет
+ `me` --- мои данные
    + `/api/v1/me` GET --- мои данные (моего пользователя)
    + `/api/v1/me/avatar` GET --- мой аватар
    + `/api/v1/me/roles` GET --- мои роли
    + `/api/v1/me/permissions` GET --- мои права доступа
    + `/api/v1/me/posts` GET --- мои объявления
    + `/api/v1/me/children` GET --- мои дети (для родителей)
    + `/api/v1/me/class` GET --- мой класс/классы моих детей (массив) --- для учеников, учителей, родителей
    + `/api/v1/me` PATCH --- обновить мои данные (моего пользователя)
        - `{email?, first_name?, last_name?, middle_name?, avatar_url?, current_password?, new_password?}`

