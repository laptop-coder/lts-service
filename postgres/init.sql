-- TODO: remove
DROP TABLE IF EXISTS 
    invite_code, post, parent_student, parent, student, 
    teacher_subject, subject, students_group, teacher, 
    room, staff, institution_administrator, user_role, 
    role_permission, permission, role, "user", 
    administrator_position, staff_position, base_entity CASCADE;

CREATE TABLE base_entity (
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        CHECK (
            created_at BETWEEN '2025-01-01'::timestamp AND CURRENT_TIMESTAMP
        ),
    -- "updated_at" is kept up to date by trigger
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        CHECK (
            updated_at BETWEEN created_at AND CURRENT_TIMESTAMP
        )
);

-- Info, related to all users (general info for all roles)
CREATE TABLE "user" (
    -- UUID4; maybe it will be used in URLs to see user's profile, e.g., or his
    -- posts
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- will be used as login
    email VARCHAR(320) NOT NULL UNIQUE
        CHECK (
            LENGTH(email) >= 5
        ),
    password VARCHAR(60) NOT NULL
        CHECK (
            LENGTH(password) = 60
        ), -- bcrypt hash
    first_name VARCHAR(100) NOT NULL
        CHECK (
            LENGTH(first_name) >= 2
        ),
    middle_name VARCHAR(100)
        CHECK (
            middle_name IS NULL OR LENGTH(middle_name) >= 2
        ),
    last_name VARCHAR(100) NOT NULL
        CHECK (
            LENGTH(last_name) >= 2
        ),
    -- Name of the avatar is ID of the user. It is located in the root of the
    -- site (public dir) and saved in the JPEG format: "/<user_id>.jpeg"
    has_avatar BOOLEAN NOT NULL DEFAULT FALSE
) INHERITS (base_entity);

-- List of users' roles (e.g., teacher, parent, student)
CREATE TABLE role (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 3
        )
) INHERITS (base_entity);

-- List of permissions (e.g., can_view_posts, can_generate_invite_codes)
CREATE TABLE permission (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 6
        )
) INHERITS (base_entity);

-- The linking table for assigning permissions to roles
CREATE TABLE role_permission (
    role_id BIGINT NOT NULL REFERENCES role (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    permission_id BIGINT NOT NULL REFERENCES permission (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    PRIMARY KEY (role_id, permission_id) -- many-to-many (role-permission)
) INHERITS (base_entity);

-- The linking table for assigning roles to users
CREATE TABLE user_role (
    user_id UUID NOT NULL REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    role_id BIGINT NOT NULL REFERENCES role (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    PRIMARY KEY (user_id, role_id) -- many-to-one (user-role)
) INHERITS (base_entity);

-- List of institution administrator positions (headmaster, e.g.)
CREATE TABLE administrator_position (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 4
        )
) INHERITS (base_entity);

-- List of staff positions (cleaner, security, etc.)
CREATE TABLE staff_position (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 4
        )
) INHERITS (base_entity);

-- Info, related only to institution administrators
CREATE TABLE institution_administrator (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one
    -- many-to-one (institution_administrator-position)
    -- can't remove position if there are at least one person with it
    position_id BIGINT NOT NULL REFERENCES administrator_position (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- Info, related only to staff
CREATE TABLE staff (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one
    -- many-to-one (staff-position)
    position_id BIGINT NOT NULL REFERENCES staff_position (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- List of rooms (cabinets, dining room, etc.)
CREATE TABLE room (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 1
        )
) INHERITS (base_entity);

-- Info, related only to teachers
CREATE TABLE teacher (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one
    -- Can't remove room if there are at least one teacher, assigned to it. To 
    -- remove the room you need to reassign the teacher to other classroom at
    -- first (actual for schools)
    classroom_id BIGINT REFERENCES room (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT -- one-to-one (teacher-room)
) INHERITS (base_entity);
    
CREATE TABLE students_group (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 1
        )
    -- Set group_advisor_id null in case of removing the user (i.e. the advisor)
    group_advisor_id UUID REFERENCES "user" (id)
        ON DELETE SET NULL
        ON UPDATE RESTRICT -- one-to-one (students_group-"user")
) INHERITS (base_entity);

-- List of subjects (e.g., "Русский язык", "Литература")
CREATE TABLE subject (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 3
        )
) INHERITS (base_entity);

-- Linking table for assigning subjects to teachers
CREATE TABLE teacher_subject (
    teacher_id UUID NOT NULL REFERENCES teacher (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    subject_id BIGINT NOT NULL REFERENCES subject (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    PRIMARY KEY (teacher_id, subject_id) -- many-to-many (teacher-subject)
) INHERITS (base_entity);

-- Info, related only to students
CREATE TABLE student (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one (student-user)
    -- many-to-one (student-students_group)
    -- Can't remove students_group if there are at least one student in it. To 
    -- remove the group you need to reassign all students to other group at
    -- first
    "group_id" BIGINT NOT NULL REFERENCES students_group (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- Info, related only to parents
CREATE TABLE parent (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY -- one-to-one (parent-user)
) INHERITS (base_entity);

-- Linking table for assigning students to their parents
CREATE TABLE parent_student (
    parent_id UUID NOT NULL REFERENCES parent (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    student_id UUID NOT NULL REFERENCES student (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT,
    PRIMARY KEY (parent_id, student_id) -- many-to-many (parent-student)
) INHERITS (base_entity);

-- Info about lost things in the format of posts
CREATE TABLE post (
    -- it will be used in URLs to see the status of the post
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Removing of the user will cause removing all of his posts
    author_id UUID NOT NULL REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT, -- many-to-one (post-author)
    name VARCHAR(50) NOT NULL
        CHECK (
            LENGTH(name) >= 2
        ),
    description VARCHAR(1000),
    -- was the post verified by moderator (service administrator)? (true/false)
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    -- was the thing found, i.e. returned to owner? (true/false)
    thing_returned_to_owner BOOLEAN NOT NULL DEFAULT FALSE
        CHECK (
            (thing_returned_to_owner = TRUE AND verified = TRUE) OR
            thing_returned_to_owner = FALSE
        ),
    -- the logic is the same as for user's avatar
    has_photo BOOLEAN NOT NULL DEFAULT FALSE
) INHERITS (base_entity);

-- Trigger to update "updated_at" column (set current time instead of the
-- previous value) in all tables, inherited from base_entity
CREATE OR REPLACE FUNCTION refresh_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER refresh_updated_at
BEFORE UPDATE ON base_entity
FOR EACH ROW EXECUTE FUNCTION refresh_updated_at_column();

-- Trigger to limit the number of superadmin accounts
CREATE OR REPLACE FUNCTION check_superadmin_limit()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role_id = (SELECT id FROM role WHERE name = 'Суперадминистратор') THEN
        IF (SELECT COUNT(*) FROM user_role WHERE role_id = NEW.role_id) >= 1 THEN
            RAISE EXCEPTION 'There can be only one superadministrator';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER superadmin_limit
BEFORE INSERT ON user_role
FOR EACH ROW EXECUTE FUNCTION check_superadmin_limit();


INSERT INTO role (name) VALUES
    ('Суперадминистратор'),
    ('Администратор сервиса'),
    ('Администрация ОУ'),
    ('Сотрудник'),
    ('Преподаватель'),
    ('Родитель'),
    ('Обучающийся');

