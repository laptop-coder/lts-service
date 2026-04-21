import { createSignal, For, Show, onMount } from "solid-js";
import { useParams } from "@solidjs/router";
import {
  usePermissions,
  PERMISSIONS,
  ROLES,
  ROLES_TO_DISPLAY,
} from "../lib/permissions";
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

const PublicProfile = () => {
  const params = useParams();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(true);
  const { hasPermission } = usePermissions();
  // Special fields
  const [teacherClassroom, setTeacherClassroom] = createSignal<Room | null>(
    null,
  );
  const [teacherSubjects, setTeacherSubjects] = createSignal<Subject[]>([]);
  const [teacherStudentGroups, setTeacherStudentGroups] = createSignal<
    StudentGroup[]
  >([]);
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
  const [studentParents, setStudentParents] = createSignal<Parent[]>([]);
  const [studentParentsUsers, setStudentParentsUsers] = createSignal<User[]>(
    [],
  );

  const { hasRole } = usePermissions();
  const [user, setUser] = createSignal<User | null>(null);

  const loadUser = async () => {
    try {
      const data = await api.get<{ user: User }>(`/users/${params.id}`);
      setUser(data.user);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось загрузить пользователя");
    } finally {
      setLoading(false);
    }
  };

  onMount(async () => {
    await loadUser();
    if (!user()) return;

    // Get special fields data (depends on roles)
    if (user()!.roles.find((r) => r.name === ROLES.INSTITUTION_ADMINISTRATOR)) {
      // institution administrator
      const institutionAdministratorData = await api.get<{
        institutionAdministrator: InstitutionAdministrator;
      }>(`/institution_administrators/${user()!.id}`);
      setInstitutionAdministratorPosition(
        institutionAdministratorData.institutionAdministrator.position || null,
      );
    }
    if (user()!.roles.find((r) => r.name === ROLES.STAFF)) {
      // staff
      const staffData = await api.get<{ staff: Staff }>(`/staff/${user()!.id}`);
      setStaffPosition(staffData.staff.position || null);
    }
    if (user()!.roles.find((r) => r.name === ROLES.TEACHER)) {
      // teacher
      const teacherData = await api.get<{ teacher: Teacher }>(
        `/teachers/${user()!.id}`,
      );
      setTeacherClassroom(teacherData.teacher.classroom || null);
      setTeacherSubjects(teacherData.teacher.subjects || []);
      setTeacherStudentGroups(teacherData.teacher.studentGroups || []);
    }
    if (user()!.roles.find((r) => r.name === ROLES.PARENT)) {
      if (hasRole(ROLES.ADMIN)) {
        // parent
        await loadParentStudents();
      }
    }
    if (user()!.roles.find((r) => r.name === ROLES.STUDENT)) {
      // student
      const studentData = await api.get<{ student: Student }>(
        `/students/${user()!.id}`,
      );
      setStudentGroup(studentData.student.studentGroup || null);
      setStudentParents(studentData.student.parents || []);

      if (hasRole(ROLES.ADMIN)) {
        // load parents data
        const studentParentsPromises = studentParents().map((parent) =>
          api.get<{ user: User }>(`/users/${parent.userId}`),
        );
        const studentParentsResponses = await Promise.all(
          studentParentsPromises,
        );
        setStudentParentsUsers(studentParentsResponses.map((r) => r.user));
      }
    }
  });

  const loadParentStudents = async () => {
    if (!user()) return;

    try {
      const data = await api.get<{ students: Student[] }>(
        `/parents/${user()!.id}/students`,
      );
      setParentStudents(data.students);

      // load students data
      const parentStudentsPromises = parentStudents().map((student) =>
        api.get<{ user: User }>(`/users/${student.userId}`),
      );
      const parentStudentsResponses = await Promise.all(parentStudentsPromises);
      setParentStudentsUsers(parentStudentsResponses.map((r) => r.user));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка загрузки учеников");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {hasPermission(PERMISSIONS.USER_READ_OWN) && (
        <div class="max-w-4xl mx-auto space-y-6 p-4">
          <h1 class="text-2xl font-bold text-center text-gray-800">
            Профиль пользователя
          </h1>

          <Show when={!user()}>
            <div class="text-center py-8 text-gray-500">Загрузка...</div>
          </Show>

          <Show when={error()}>
            <div class="bg-red-100 text-red-700 p-4 rounded-xl">{error()}</div>
          </Show>
          <Show when={user() && !error()}>
            <div class="bg-white rounded-2xl shadow-lg p-6">
              <div class="flex flex-col md:flex-row gap-6 items-center md:items-start">
                <div class="relative group w-32 h-32 rounded-full">
                  <img
                    class="w-32 h-32 rounded-full object-cover border-4 border-blue-100"
                    src={`/storage/storage/avatars/${user()!.hasAvatar ? user()!.id : "default"}.jpeg`}
                    alt="Фото профиля"
                  />
                </div>
                <div class="flex-1 text-center md:text-left">
                  <h2 class="text-2xl font-bold text-gray-800">
                    {user()!.lastName} {user()!.firstName} {user()?.middleName}
                  </h2>
                  <p class="text-gray-500 mt-1">{user()!.email}</p>
                  <div class="flex flex-wrap gap-2 mt-3">
                    <div class="flex flex-wrap gap-1">
                      <For each={user()!.roles}>
                        {(ur) => (
                          <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">
                            {
                              ROLES_TO_DISPLAY.find((r) => r.id === ur.id)!
                                .displayName
                            }
                          </span>
                        )}
                      </For>
                    </div>
                  </div>
                </div>
                <div class="text-sm text-gray-500">
                  <Show when={hasRole(ROLES.ADMIN)}>
                    <p>ID: {user()!.id}</p>
                    <p>Аккаунт создан: {formatDate(user()!.createdAt)}</p>
                  </Show>
                </div>
              </div>
            </div>

            <Show
              when={
                hasRole(ROLES.ADMIN) &&
                user()!.roles.find((r) => r.name === ROLES.PARENT)
              }
            >
              <Show when={parentStudentsUsers().length > 0}>
                <div class="bg-white rounded-2xl shadow-lg p-6">
                  <h2 class="text-xl font-bold text-gray-800 mb-4">Дети</h2>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <For each={parentStudentsUsers()}>
                      {(user) => (
                        <div class="border rounded-xl p-4 hover:shadow-md transition relative">
                          <div class="flex items-center gap-3">
                            <img
                              class="w-12 h-12 rounded-full object-cover"
                              src={`/storage/storage/avatars/${user.hasAvatar ? user.id : "default"}.jpeg`}
                              alt="Фото профиля"
                            />
                            <div>
                              <p class="font-semibold">
                                {user.lastName} {user.firstName}{" "}
                                {user?.middleName}
                              </p>
                              <p class="text-sm text-gray-500">{user.email}</p>

                              <div class="flex flex-wrap gap-2 mt-3 mb-3">
                                <div class="flex flex-wrap gap-1">
                                  <For
                                    each={ROLES_TO_DISPLAY.filter((rd) =>
                                      user.roles.some((ur) => ur.id === rd.id),
                                    )}
                                  >
                                    {(role) => (
                                      <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">
                                        {role.displayName}
                                      </span>
                                    )}
                                  </For>
                                </div>
                              </div>
                              <div class="text-sm text-gray-500">
                                <p>ID: {user.id}</p>
                                <p>
                                  Аккаунт создан: {formatDate(user.createdAt)}
                                </p>
                              </div>
                            </div>
                          </div>
                        </div>
                      )}
                    </For>
                  </div>
                </div>
              </Show>
            </Show>
            <Show when={user()!.roles.find((r) => r.name === ROLES.TEACHER)}>
              <div class="bg-white rounded-2xl shadow-lg p-6 space-y-4">
                <h3 class="text-lg font-semibold text-gray-700">
                  Преподаватель
                </h3>

                <div>
                  <h4 class="text-sm font-medium text-gray-500 mb-2">
                    Предметы
                  </h4>
                  <div class="flex flex-wrap gap-2">
                    <For each={teacherSubjects()}>
                      {(subject) => (
                        <span class="px-3 py-1 bg-green-100 text-green-700 text-sm rounded-full">
                          {subject.name}
                        </span>
                      )}
                    </For>
                    <Show when={teacherSubjects().length === 0}>
                      <span class="text-gray-500 text-sm">Нет предметов</span>
                    </Show>
                  </div>
                </div>

                <div>
                  <h4 class="text-sm font-medium text-gray-500 mb-2">
                    Кабинет
                  </h4>
                  <div class="flex items-center gap-3">
                    <span class="w-2 h-2 bg-green-500 rounded-full"></span>
                    <span class="text-gray-800">
                      {teacherClassroom()?.name || "Не указан"}
                    </span>
                  </div>
                </div>

                <div>
                  <h4 class="text-sm font-medium text-gray-500 mb-2">
                    Классное руководство/менторство
                  </h4>
                  <div class="flex flex-wrap gap-2">
                    <For each={teacherStudentGroups()}>
                      {(group) => (
                        <span class="px-3 py-1 bg-emerald-100 text-emerald-700 text-sm rounded-full">
                          {group.name}
                        </span>
                      )}
                    </For>
                    <Show when={teacherStudentGroups().length === 0}>
                      <span class="text-gray-500 text-sm">
                        Нет учебных групп или классов
                      </span>
                    </Show>
                  </div>
                </div>
              </div>
            </Show>
            <Show
              when={user()!.roles.find(
                (r) => r.name === ROLES.INSTITUTION_ADMINISTRATOR,
              )}
            >
              <div class="bg-white rounded-2xl shadow-lg p-6">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">
                  Должность
                </h3>
                <div class="flex items-center gap-3">
                  <span class="w-2 h-2 bg-red-500 rounded-full"></span>
                  <span class="text-gray-800">
                    {institutionAdministratorPosition()?.name}
                  </span>
                </div>
              </div>
            </Show>
            <Show when={user()!.roles.find((r) => r.name === ROLES.STAFF)}>
              <div class="bg-white rounded-2xl shadow-lg p-6">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">
                  Должность
                </h3>
                <div class="flex items-center gap-3">
                  <span class="w-2 h-2 bg-indigo-500 rounded-full"></span>
                  <span class="text-gray-800">{staffPosition()?.name}</span>
                </div>
              </div>
            </Show>
            <Show when={user()!.roles.find((r) => r.name === ROLES.STUDENT)}>
              <Show
                when={hasRole(ROLES.ADMIN) && studentParentsUsers().length > 0}
              >
                <div class="bg-white rounded-2xl shadow-lg p-6">
                  <h2 class="text-xl font-bold text-gray-800 mb-4">
                    Родители ученика
                  </h2>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <For each={studentParentsUsers()}>
                      {(user) => (
                        <div class="border rounded-xl p-4 hover:shadow-md transition">
                          <div class="flex items-center gap-3">
                            <img
                              class="w-12 h-12 rounded-full object-cover"
                              src={`/storage/storage/avatars/${user.hasAvatar ? user.id : "default"}.jpeg`}
                              alt="Фото профиля"
                            />
                            <div>
                              <p class="font-semibold">
                                {user.lastName} {user.firstName}{" "}
                                {user?.middleName}
                              </p>
                              <p class="text-sm text-gray-500">{user.email}</p>

                              <div class="flex flex-wrap gap-2 mt-3 mb-3">
                                <div class="flex flex-wrap gap-1">
                                  <For
                                    each={ROLES_TO_DISPLAY.filter((rd) =>
                                      user.roles.some((ur) => ur.id === rd.id),
                                    )}
                                  >
                                    {(role) => (
                                      <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">
                                        {role.displayName}
                                      </span>
                                    )}
                                  </For>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      )}
                    </For>
                  </div>
                </div>
              </Show>
              <div class="bg-white rounded-2xl shadow-lg p-6">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">
                  Класс/учебная группа
                </h3>
                <div class="flex items-center gap-3">
                  <span class="w-2 h-2 bg-pink-500 rounded-full"></span>
                  <span class="text-gray-800">{studentGroup()?.name}</span>
                </div>
              </div>
            </Show>
          </Show>
        </div>
      )}
    </>
  );
};

export default PublicProfile;
