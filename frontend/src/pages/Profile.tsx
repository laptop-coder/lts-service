import { createSignal, For, Show, onMount, Index } from "solid-js";
import { createStore } from "solid-js/store";
import { useNavigate } from "@solidjs/router";
import {
  usePermissions,
  PERMISSIONS,
  ROLES,
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
  const auth = useAuth();
  const navigate = useNavigate();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
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
  const [uploadingAvatar, setUploadingAvatar] = createSignal(false);

  const { user } = useAuth();
  const { hasRole } = usePermissions();

  // For profile edit
  const [editMode, setEditMode] = createSignal(false);
  const [saving, setSaving] = createSignal(false);
  // Data for <select> tags
  const [rooms, setRooms] = createSignal<Room[]>([]);
  const [subjects, setSubjects] = createSignal<Subject[]>([]);
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  const [
    institutionAdministratorPositions,
    setInstitutionAdministratorPositions,
  ] = createSignal<InstitutionAdministratorPosition[]>([]);
  // Fields for <select> tags
  const [teacherClassroomId, setTeacherClassroomId] = createSignal<
    number | null
  >(null);
  const [teacherSubjectIds, setTeacherSubjectIds] = createSignal<number[]>([]);
  const [teacherStudentGroupIds, setTeacherStudentGroupIds] = createStore<
    number[]
  >([]);
  const [studentGroupId, setStudentGroupId] = createSignal<number | null>(null);
  const [staffPositionId, setStaffPositionId] = createSignal<number | null>(
    null,
  );
  const [
    institutionAdministratorPositionId,
    setInstitutionAdministratorPositionId,
  ] = createSignal<number | null>(null);
  const [parentStudentIds, setParentStudentIds] = createStore<string[]>([]);

  const loadDataForSelect = async () => {
    const [
      roomsData,
      subjectsData,
      groupsData,
      staffPositionData,
      institutionAdministratorPositionData,
    ] = await Promise.all([
      api.get<{ rooms: Room[] }>("/rooms?limit=65535"), // uint16
      api.get<{ subjects: Subject[] }>("/subjects?limit=65535"),
      api.get<{ studentGroups: StudentGroup[] }>("/student_groups?limit=65535"),
      api.get<{ staffPositions: StaffPosition[] }>(
        "/staff/positions?limit=65535",
      ),
      api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions?limit=65535"),
    ]);
    setRooms(roomsData.rooms);
    setSubjects(subjectsData.subjects);
    setStudentGroups(groupsData.studentGroups);
    setStaffPositions(staffPositionData.staffPositions);
    setInstitutionAdministratorPositions(
      institutionAdministratorPositionData.institutionAdministratorPositions,
    );
  };

  const loadAllData = async () => {
    setLoading(true);
    await loadDataForSelect();

    // Get special fields data (depends on roles)
    if (hasRole(ROLES.INSTITUTION_ADMINISTRATOR)) {
      // institution administrator
      const institutionAdministratorData = await api.get<{
        institutionAdministrator: InstitutionAdministrator;
      }>(`/institution_administrators/${user()!.id}`);
      setInstitutionAdministratorPosition(
        institutionAdministratorData.institutionAdministrator.position || null,
      );
      setInstitutionAdministratorPositionId(
        institutionAdministratorData.institutionAdministrator.position.id ||
          null,
      );
    }
    if (hasRole(ROLES.STAFF)) {
      // staff
      const staffData = await api.get<{ staff: Staff }>(`/staff/${user()!.id}`);
      setStaffPosition(staffData.staff.position || null);
      setStaffPositionId(staffData.staff.position.id || null);
    }
    if (hasRole(ROLES.TEACHER)) {
      // teacher
      const teacherData = await api.get<{ teacher: Teacher }>(
        `/teachers/${user()!.id}`,
      );
      setTeacherClassroom(teacherData.teacher.classroom || null);
      setTeacherSubjects(teacherData.teacher.subjects || []);
      setTeacherStudentGroups(teacherData.teacher.studentGroups || []);
      setTeacherClassroomId(teacherData.teacher.classroom?.id || null);
      setTeacherSubjectIds(teacherData.teacher.subjects.map((s) => s.id) || []);
      setTeacherStudentGroupIds(
        teacherData.teacher.studentGroups?.map((g) => g.id) || [],
      );
    }
    if (hasRole(ROLES.PARENT)) {
      // parent
      await loadParentStudents();
    }
    if (hasRole(ROLES.STUDENT)) {
      // student
      const studentData = await api.get<{ student: Student }>(
        `/students/${user()!.id}`,
      );
      setStudentGroup(studentData.student.studentGroup || null);
      setStudentParents(studentData.student.parents || []);
      setStudentGroupId(studentData.student.studentGroup.id || null);
      setStudentParents(studentData.student.parents || []);

      // load parents data
      const studentParentsPromises = studentParents().map((parent) =>
        api.get<{ user: User }>(`/users/${parent.userId}`),
      );
      const studentParentsResponses = await Promise.all(studentParentsPromises);
      setStudentParentsUsers(studentParentsResponses.map((r) => r.user));
    }
    setLoading(false);
  };

  onMount(async () => {
    if (!user()) return;
    await loadAllData();
  });

  const handleLogout = async () => {
    await auth.logout();
    navigate("/login");
  };

  const loadParentStudents = async () => {
    try {
      const data = await api.get<{ students: Student[] }>(
        "/parents/me/students",
      );
      setParentStudents(data.students);
      setParentStudentIds(data.students.map((s) => s.userId));

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

  // Parent students
  const addStudentId = () => {
    setSaving(false);
    setError("");
    setParentStudentIds([...parentStudentIds, ""]);
  };
  const updateStudentId = (index: number, value: string) => {
    setSaving(false);
    setError("");
    setParentStudentIds(index, value);
  };
  const removeStudentId = (index: number) => {
    setSaving(false);
    setError("");
    setParentStudentIds(parentStudentIds.filter((_, i) => i !== index));
  };

  // Student groups where teacher is the advisor
  const addTeacherStudentGroupId = () => {
    setTeacherStudentGroupIds([
      ...teacherStudentGroupIds.map((id) => Number(id)),
      0,
    ]);
  };
  const updateTeacherStudentGroupId = (index: number, value: number) => {
    setTeacherStudentGroupIds(index, value);
  };
  const removeTeacherStudentGroupId = (index: number) => {
    setTeacherStudentGroupIds(
      teacherStudentGroupIds.filter((_, i) => i !== index),
    );
  };

  const handleAvatarUpload = async (e: Event) => {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    setUploadingAvatar(true);
    const formData = new FormData();
    formData.append("avatar", file);

    try {
      await api.put("/users/me/avatar", formData);
      // Update user data
      const userData = await api.get<{ user: User }>("/users/me");
      auth.user()!.hasAvatar = userData.user.hasAvatar;
      // Reload page
      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось обновить аватар",
      );
    } finally {
      setUploadingAvatar(false);
    }
  };

  const cancelEdit = async () => {
    setEditMode(false);
    await loadAllData();
  };

  const saveProfile = async () => {
    const formData = new URLSearchParams();

    if (hasRole(ROLES.INSTITUTION_ADMINISTRATOR)) {
      if (institutionAdministratorPositionId())
        formData.append(
          "institutionAdministratorPositionId",
          String(institutionAdministratorPositionId()),
        );
    }

    if (hasRole(ROLES.STAFF)) {
      if (staffPositionId())
        formData.append("staffPositionId", String(staffPositionId()));
    }

    if (hasRole(ROLES.TEACHER)) {
      if (teacherClassroomId())
        formData.append("teacherClassroomId", String(teacherClassroomId()));
      teacherSubjectIds().forEach((id) =>
        formData.append("teacherSubjectId", String(id)),
      );
      teacherStudentGroupIds.forEach((id) =>
        formData.append("teacherStudentGroupId", String(id)),
      );
    }

    if (hasRole(ROLES.PARENT)) {
      parentStudentIds.forEach((studentId) => {
        if (studentId.trim()) {
          formData.append("parentStudentId", studentId);
        }
      });
    }

    if (hasRole(ROLES.STUDENT)) {
      if (studentGroup())
        formData.append("studentGroupId", String(studentGroupId()));
    }

    await api.put("/users/me/extensions", formData);
    setEditMode(false);
    await loadAllData();
  };

  return (
    <>
      {hasPermission(PERMISSIONS.USER_READ_OWN) && (
        <div class="max-w-4xl mx-auto space-y-6 p-4">
          <h1 class="text-2xl font-bold text-center text-gray-800">
            Мой профиль
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
                  <label class="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full opacity-0 group-hover:opacity-100 transition cursor-pointer">
                    <svg
                      class="w-8 h-8 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                      />
                    </svg>
                    <input
                      type="file"
                      accept="image/jpeg,image/png,image/webp,image/gif"
                      onChange={handleAvatarUpload}
                      class="hidden"
                      disabled={uploadingAvatar()}
                    />
                  </label>
                  <Show when={uploadingAvatar()}>
                    <div class="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full">
                      <div class="text-white text-sm">Загрузка...</div>
                    </div>
                  </Show>
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
                  <p>ID: {user()!.id}</p>
                  <p>Аккаунт создан: {formatDate(user()!.createdAt)}</p>
                  <div class="flex gap-1 flex-col">
                    <p>
                      <button
                        onClick={handleLogout}
                        class="px-3 py-1.5 bg-red-700 text-white rounded-lg hover:bg-red-800 transition mt-5 cursor-pointer"
                      >
                        Выйти
                      </button>
                    </p>
                    <div class="flex gap-3 flex-row">
                      <Show when={!editMode()}>
                        <button
                          onClick={() => setEditMode(true)}
                          class="px-3 py-1.5 bg-blue-700 text-white rounded-lg hover:bg-blue-800 transition mt-5 cursor-pointer"
                        >
                          Редактировать
                        </button>
                      </Show>
                      <Show when={editMode()}>
                        <button
                          onClick={cancelEdit}
                          class="px-3 py-1.5 bg-red-700 text-white rounded-lg hover:bg-red-800 transition mt-5 cursor-pointer"
                        >
                          Отмена
                        </button>
                        <button
                          onClick={saveProfile}
                          class="px-3 py-1.5 bg-green-700 text-white rounded-lg hover:bg-green-800 transition mt-5 cursor-pointer"
                        >
                          Сохранить
                        </button>
                      </Show>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <Show when={hasRole(ROLES.PARENT)}>
              <Show when={parentStudentsUsers().length > 0 || editMode()}>
                <div class="bg-white rounded-2xl shadow-lg p-6">
                  <h2 class="text-xl font-bold text-gray-800">Мои дети</h2>
                  <Show when={editMode()}>
                    <div class="space-y-3 border-t border-gray-100 pt-4">
                      <h3 class="font-medium text-gray-800">
                        Привязка учеников
                      </h3>

                      <Index each={parentStudentIds}>
                        {(studentId, index) => (
                          <div class="flex gap-2">
                            <input
                              disabled={saving()}
                              type="text"
                              value={studentId()}
                              onInput={(e) => {
                                setSaving(false);
                                setError("");
                                updateStudentId(index, e.target.value);
                              }}
                              placeholder={`ID ученика ${index + 1}`}
                              class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50 disabled:cursor-not-allowed"
                            />
                            <button
                              disabled={saving()}
                              type="button"
                              onClick={() => removeStudentId(index)}
                              class="px-4 py-2 bg-red-700 text-white rounded-xl hover:bg-red-800 transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                            >
                              Удалить
                            </button>
                          </div>
                        )}
                      </Index>

                      <button
                        disabled={saving()}
                        type="button"
                        onClick={addStudentId}
                        class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                      >
                        + Добавить ученика
                      </button>
                    </div>
                  </Show>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-5 mb-5">
                    <For each={parentStudentsUsers()}>
                      {(user) => (
                        <div
                          class={`border rounded-xl p-4 hover:shadow-md transition relative ${!parentStudentIds.includes(user.id) ? "opacity-50" : ""}`}
                        >
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
            <Show when={hasRole(ROLES.TEACHER)}>
              <div class="bg-white rounded-2xl shadow-lg p-6 space-y-4">
                <h3 class="text-lg font-semibold text-gray-700">
                  Преподаватель
                </h3>

                <Show when={!editMode()}>
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
                </Show>
                <Show when={editMode()}>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Классный кабинет *
                    </label>
                    <select
                      value={teacherClassroomId() || ""}
                      onChange={(e) =>
                        setTeacherClassroomId(Number(e.currentTarget.value))
                      }
                      class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <option value="">Выберите кабинет</option>
                      <For each={rooms()}>
                        {(room) => <option value={room.id}>{room.name}</option>}
                      </For>
                    </select>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Предметы *
                    </label>
                    <div class="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto p-2 border border-gray-200 rounded-xl">
                      <For each={subjects()}>
                        {(subject) => (
                          <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded-lg cursor-pointer transition">
                            <input
                              type="checkbox"
                              checked={teacherSubjectIds().includes(subject.id)}
                              onChange={() => {
                                if (teacherSubjectIds().includes(subject.id)) {
                                  setTeacherSubjectIds((prev) =>
                                    prev.filter((id) => id !== subject.id),
                                  );
                                } else {
                                  setTeacherSubjectIds((prev) => [
                                    ...prev,
                                    subject.id,
                                  ]);
                                }
                              }}
                              class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                            />
                            <span class="text-gray-700 text-sm">
                              {subject.name}
                            </span>
                          </label>
                        )}
                      </For>
                    </div>
                  </div>
                  <div class="space-y-3 border-t border-gray-200 pt-4">
                    <h3 class="font-medium text-gray-800">
                      Классное руководство/менторство
                    </h3>

                    <Index each={teacherStudentGroupIds}>
                      {(groupId, index) => (
                        <div class="flex gap-2">
                          <select
                            value={groupId() || ""}
                            onChange={(e) =>
                              updateTeacherStudentGroupId(
                                index,
                                Number(e.target.value),
                              )
                            }
                            class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            <option value="">Выберите группу</option>
                            <For
                              each={studentGroups().filter(
                                (group) =>
                                  !teacherStudentGroupIds
                                    .filter((_, i) => i !== index)
                                    .includes(group.id),
                              )}
                            >
                              {(group) => (
                                <option value={group.id}>{group.name}</option>
                              )}
                            </For>
                          </select>
                          <button
                            type="button"
                            onClick={() => removeTeacherStudentGroupId(index)}
                            class="px-4 py-2 bg-red-700 text-white rounded-xl hover:bg-red-800 transition cursor-pointer"
                          >
                            Удалить
                          </button>
                        </div>
                      )}
                    </Index>

                    <button
                      type="button"
                      onClick={addTeacherStudentGroupId}
                      class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
                    >
                      + Добавить группу
                    </button>
                  </div>
                </Show>
              </div>
            </Show>
            <Show when={hasRole(ROLES.INSTITUTION_ADMINISTRATOR)}>
              <div class="bg-white rounded-2xl shadow-lg p-6">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">
                  Должность
                </h3>
                <Show when={!editMode()}>
                  <div class="flex items-center gap-3">
                    <span class="w-2 h-2 bg-red-500 rounded-full"></span>
                    <span class="text-gray-800">
                      {institutionAdministratorPosition()?.name}
                    </span>
                  </div>
                </Show>
                <Show when={editMode()}>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Должность администрации ОУ *
                    </label>
                    <select
                      disabled={saving()}
                      value={institutionAdministratorPositionId() || ""}
                      onChange={(e) => {
                        setSaving(false);
                        setError("");
                        setInstitutionAdministratorPositionId(
                          Number(e.currentTarget.value),
                        );
                      }}
                      class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                      <option value="">Выберите должность</option>
                      <For each={institutionAdministratorPositions()}>
                        {(position) => (
                          <option value={position.id}>{position.name}</option>
                        )}
                      </For>
                    </select>
                  </div>
                </Show>
              </div>
            </Show>
            <Show when={hasRole(ROLES.STAFF)}>
              <div class="bg-white rounded-2xl shadow-lg p-6">
                <h3 class="text-lg font-semibold text-gray-700 mb-4">
                  Должность
                </h3>
                <Show when={!editMode()}>
                  <div class="flex items-center gap-3">
                    <span class="w-2 h-2 bg-indigo-500 rounded-full"></span>
                    <span class="text-gray-800">{staffPosition()?.name}</span>
                  </div>
                </Show>
                <Show when={editMode()}>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Должность сотрудника ОУ *
                    </label>
                    <select
                      disabled={saving()}
                      value={staffPositionId() || ""}
                      onChange={(e) => {
                        setSaving(false);
                        setError("");
                        setStaffPositionId(Number(e.currentTarget.value));
                      }}
                      class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                      <option value="">Выберите должность</option>
                      <For each={staffPositions()}>
                        {(position) => (
                          <option value={position.id}>{position.name}</option>
                        )}
                      </For>
                    </select>
                  </div>
                </Show>
              </div>
            </Show>
            <Show when={hasRole(ROLES.STUDENT)}>
              <Show when={studentParentsUsers().length > 0}>
                <div class="bg-white rounded-2xl shadow-lg p-6">
                  <h2 class="text-xl font-bold text-gray-800 mb-4">
                    Мои родители
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
                <Show when={!editMode()}>
                  <div class="flex items-center gap-3">
                    <span class="w-2 h-2 bg-pink-500 rounded-full"></span>
                    <span class="text-gray-800">{studentGroup()?.name}</span>
                  </div>
                </Show>
                <Show when={editMode()}>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Учебная группа *
                    </label>
                    <select
                      disabled={saving()}
                      value={studentGroupId() || ""}
                      onChange={(e) =>
                        setStudentGroupId(Number(e.currentTarget.value))
                      }
                      class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                      <option value="">Выберите группу</option>
                      <For each={studentGroups()}>
                        {(group) => (
                          <option value={group.id}>{group.name}</option>
                        )}
                      </For>
                    </select>
                  </div>
                </Show>
              </div>
            </Show>
          </Show>
        </div>
      )}
    </>
  );
};

export default Profile;
