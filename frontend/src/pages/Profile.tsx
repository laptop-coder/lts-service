import { createSignal, For, Show, onMount } from "solid-js";
import {
  usePermissions,
  PERMISSIONS,
  ROLES_TO_DISPLAY,
} from "../lib/permissions";
import { useAuth } from "../lib/auth";
import { api } from "../lib/api";
import { formatDate } from "../lib/utils";
import type {
  User,
  Subject,
  Room,
  Parent,
  Student,
  Teacher,
  InstitutionAdministrator,
  Staff,
  StudentGroup,
  StaffPosition,
  InstitutionAdministratorPosition,
} from "../lib/types";

const Profile = () => {
  const [error, setError] = createSignal("");
  const { hasPermission } = usePermissions();
  // Special fields
  const [teacherClassroom, setTeacherClassroom] = createSignal<Room | null>(
    null,
  );
  const [teacherSubjects, setTeacherSubjects] = createSignal<Subject[]>([]);
  const [studentGroup, setStudentGroup] = createSignal<StudentGroup | null>(
    null,
  );
  const [staffPosition, setStaffPosition] = createSignal<StaffPosition | null>(
    null,
  );
  const [
    institutionAdministratorPosition,
    setInstitutionAdministratorPosition,
  ] = createSignal<InstitutionAdministratorPosition | null>(null);
  const [parentStudents, setParentStudents] = createSignal<Student[]>([]);
  const [parentStudentsUsers, setParentStudentsUsers] = createSignal<User[]>(
    [],
  );
  const [loadingParentStudentsUsers, setLoadingParentStudentsUsers] =
    createSignal(false);

  const { user } = useAuth();

  onMount(async () => {
    if (!user()) return;

    // Get special fields data (depends on roles)
    if (user()!.roles.some((r) => r.id === 3)) {
      // institution administrator
      const institutionAdministratorData = await api.get<{
        institutionAdministrator: InstitutionAdministrator;
      }>(`/institution_administrators/${user()!.id}`);
      setInstitutionAdministratorPosition(
        institutionAdministratorData.institutionAdministrator.position || null,
      );
    }
    if (user()!.roles.some((r) => r.id === 4)) {
      // staff
      const staffData = await api.get<{ staff: Staff }>(`/staff/${user()!.id}`);
      setStaffPosition(staffData.staff.position || null);
    }
    if (user()!.roles.some((r) => r.id === 5)) {
      // teacher
      const teacherData = await api.get<{ teacher: Teacher }>(
        `/teachers/${user()!.id}`,
      );
      setTeacherClassroom(teacherData.teacher.classroom || null);
      setTeacherSubjects(teacherData.teacher.subjects || []);
    }
    if (user()!.roles.some((r) => r.id === 6)) {
      // parent
      const parentData = await api.get<{ parent: Parent }>(
        `/parents/${user()!.id}`,
      );
      setParentStudents(parentData.parent.students || []);

      // load students data
      setLoadingParentStudentsUsers(true);
      const parentStudentsPromises = parentStudents().map((student) =>
        api.get<{ user: User }>(`/users/${student.userId}`),
      );
      const parentStudentsResponses = await Promise.all(parentStudentsPromises);
      setParentStudentsUsers(parentStudentsResponses.map((r) => r.user));
      setLoadingParentStudentsUsers(false);
    }
    if (user()!.roles.some((r) => r.id === 7)) {
      // student
      const studentData = await api.get<{ student: Student }>(
        `/students/${user()!.id}`,
      );
      setStudentGroup(studentData.student.studentGroup || null);
    }
  });

  return (
    <>
      {hasPermission(PERMISSIONS.USER_READ_OWN) && (
        <div class="max-w-4xl mx-auto space-y-6">
          <h1 class="text-2xl font-bold text-center">Профиль</h1>

          <Show when={!user()}>
            <div class="text-center py-8">Загрузка...</div>
          </Show>

          <Show when={error()}>
            <div class="bg-red-100 text-red-700 p-4 rounded-lg">{error()}</div>
          </Show>

          <Show when={user() && !error()}>
            <ul>
              <img
                class="w-10 h-10"
                src={`/storage/storage/avatars/${user()!.hasAvatar ? user()!.id : "default"}.jpeg`}
                alt="Фото профиля"
              />
              <li>ID пользователя: {user()!.id}</li>
              <li>Email: {user()!.email}</li>
              <li>Имя: {user()!.firstName}</li>
              <li>Фамилия: {user()!.lastName}</li>
              {user()!.middleName && <li>Отчество: {user()!.middleName}</li>}
              <li>Аккаунт создан: {formatDate(user()!.createdAt)}</li>
              <li>
                Роли:
                <ul>
                  <For
                    each={ROLES_TO_DISPLAY.filter((rd) =>
                      user()!
                        .roles.map((r) => r.id)
                        .includes(rd.id),
                    )}
                  >
                    {(role) => <li>{role.displayName}</li>}
                  </For>
                </ul>
              </li>
            </ul>

            <Show when={user()!.roles.some((r) => r.id === 3)}>
              {institutionAdministratorPosition()?.name}
            </Show>
            <Show when={user()!.roles.some((r) => r.id === 4)}>
              {staffPosition()?.name}
            </Show>
            <Show when={user()!.roles.some((r) => r.id === 5)}>
              {teacherClassroom()?.name}
              <For each={teacherSubjects()}>{(subject) => subject.name}</For>
            </Show>
            <Show when={user()!.roles.some((r) => r.id === 6)}>
              <h2 class="text-xl font-bold text-center">Дети</h2>
              <For each={parentStudentsUsers()}>
                {(user) => (
                  <ul>
                    <img
                      class="w-10 h-10"
                      src={`/storage/storage/avatars/${user.hasAvatar ? user.id : "default"}.jpeg`}
                      alt="Фото профиля"
                    />
                    <li>ID: {user.id}</li>
                    <li>
                      {user.lastName} {user.firstName} {user?.middleName}
                    </li>
                    <li>Зарегистрирован: {formatDate(user.createdAt)}</li>
                    <li>Email: {user.email}</li>
                    <li>
                      Роли:{" "}
                      <ul>
                        <For
                          each={ROLES_TO_DISPLAY.filter((rd) =>
                            user.roles.some((ur) => ur.id === rd.id),
                          )}
                        >
                          {(role) => role.displayName}
                        </For>
                      </ul>
                    </li>
                  </ul>
                )}
              </For>
            </Show>
            <Show when={user()!.roles.some((r) => r.id === 7)}>
              {studentGroup()?.name}
            </Show>
          </Show>
        </div>
      )}
    </>
  );
};

export default Profile;
