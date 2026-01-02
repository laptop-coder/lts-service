-- TODO: remove
DROP TABLE IF EXISTS 
    invite_code, post, parent_student, parent, student, 
    teacher_subject, subject, students_group, teacher, 
    room, staff, school_administrator, user_role, 
    role_permission, permission, role, "user", 
    administrator_position, staff_position, superadmin, base_entity CASCADE;

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

-- There can be only one superadmin account
CREATE TABLE superadmin (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login VARCHAR(16) NOT NULL UNIQUE
        CHECK (
            LENGTH(login) >= 5
        ),
    -- bcrypt hash always has a length of 60 bytes
    password VARCHAR(60) NOT NULL
        CHECK (
            LENGTH(password) = 60
        )
) INHERITS (base_entity);

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
    -- <UUID4>.jpeg or <UUID4>.jpg, where <UUID4> is the user's id (it has the
    -- length of 36 bytes)
    avatar_url VARCHAR(41)
        CHECK (
            (
                LENGTH(avatar_url) BETWEEN 40 AND 41
            )
            AND (
                RIGHT(avatar_url, 5) = '.jpeg'
                OR RIGHT(avatar_url, 4) = '.jpg'
            )
        )
) INHERITS (base_entity);

-- List of users' roles (e.g., teacher, parent, student)
CREATE TABLE role (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 3
        )
) INHERITS (base_entity);

-- List of permissions (e.g., can_view_posts, can_generate_invite_codes)
CREATE TABLE permission (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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

-- List of school administrator positions (headmaster, e.g.)
CREATE TABLE administrator_position (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 4
        )
) INHERITS (base_entity);

-- List of staff positions (cleaner, security, etc.)
CREATE TABLE staff_position (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 4
        )
) INHERITS (base_entity);

-- Info, related only to school administrators
CREATE TABLE school_administrator (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one
    -- many-to-one (school_administrator-position)
    -- can't remove position if there are at least one person with it
    position BIGINT NOT NULL REFERENCES administrator_position (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- Info, related only to staff
CREATE TABLE staff (
    id UUID REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT PRIMARY KEY, -- one-to-one
    -- many-to-one (staff-position)
    position BIGINT NOT NULL REFERENCES staff_position (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- List of rooms (cabinets, dining room, etc.)
CREATE TABLE room (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
    -- first
    classroom BIGINT NOT NULL REFERENCES room (id)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT -- one-to-one (teacher-room)
) INHERITS (base_entity);
    
CREATE TABLE students_group (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
        CHECK (
            LENGTH(name) >= 1
        ),
    -- Set classroom_teacher null in case of removing the teacher
    classroom_teacher UUID REFERENCES teacher (id)
        ON DELETE SET NULL
        ON UPDATE RESTRICT -- one-to-one (group-teacher)
) INHERITS (base_entity);

-- List of subjects (e.g., "Русский язык", "Литература")
CREATE TABLE subject (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
    "group" BIGINT NOT NULL REFERENCES students_group (id)
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
    author UUID NOT NULL REFERENCES "user" (id)
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
) INHERITS (base_entity);

-- List of the one-time invite codes (they are removed after use). To generate
-- the code you just need to insert role_id into this table and get generated
-- primary key.
CREATE TABLE invite_code (
    -- this codes (UUIDs) will be used in invite-links
    code UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Removing of the role will cause removing all of the invite codes with
    -- this role
    role_id BIGINT NOT NULL REFERENCES role (id)
        ON DELETE CASCADE
        ON UPDATE RESTRICT
) INHERITS (base_entity);

-- Function and trigger to update "updated_at" column (set current time instead
-- of the previous value) in all tables, inherited from base_entity
CREATE OR REPLACE FUNCTION refresh_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER refresh_updated_at
BEFORE UPDATE ON base_entity
FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at_column();


INSERT INTO role (name) VALUES
    ('Администратор сервиса'),
    ('Администрация школы'),
    ('Сотрудник'),
    ('Учитель'),
    ('Родитель'),
    ('Ученик');

