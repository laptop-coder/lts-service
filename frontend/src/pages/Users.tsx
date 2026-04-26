import {
  createSignal,
  createEffect,
  Show,
  For,
  Index,
  onMount,
  onCleanup,
} from "solid-js";
import { createStore } from "solid-js/store";
import { api } from "../lib/api";
import { useAuth } from "../lib/auth";
import {
  PERMISSIONS,
  ROLES,
  usePermissions,
  ROLES_TO_DISPLAY,
} from "../lib/permissions";
import type {
  User,
  Room,
  Subject,
  StudentGroup,
  StaffPosition,
  InstitutionAdministratorPosition,
  Parent,
  Student,
  Teacher,
  InstitutionAdministrator,
  Staff,
} from "../lib/types";
import { A } from "@solidjs/router";
import Pagination from "../components/Pagination";
import { Users as UsersIcon } from "lucide-solid";
import { Trash, Plus } from "lucide-solid";

const Users = () => {
  const auth = useAuth();
  const [users, setUsers] = createSignal<User[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [selectedUser, setSelectedUser] = createSignal<User | null>(null);
  const [selectedRoles, setSelectedRoles] = createSignal<number[]>([]);
  const [saving, setSaving] = createSignal(false);
  const [error, setError] = createSignal("");
  // Special fields
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
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);

  createEffect(() => {
    page();
    loadUsers();
  });

  // Data for special fields
  // rooms:
  const [rooms, setRooms] = createSignal<Room[]>([]);
  // subjects:
  const [subjects, setSubjects] = createSignal<Subject[]>([]);
  // student groups:
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  // staff positions:
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  // institution administrators positions:
  const [
    institutionAdministratorPositions,
    setInstitutionAdministratorPositions,
  ] = createSignal<InstitutionAdministratorPosition[]>([]);

  const { hasPermission, hasRole } = usePermissions();

  const limit = 20;

  const loadUsers = async () => {
    try {
      const data = await api.get<{ users: User[] }>(
        `/users?limit=${limit}&offset=${page() * limit}`,
      );
      setHasMore(data.users.length > limit);
      setUsers(data.users.slice(0, limit));
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки пользователей",
      );
    } finally {
      setLoading(false);
    }
  };

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

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Escape" && selectedUser()) {
      closeModal();
    }
  };

  onMount(() => {
    loadUsers();
    loadDataForSelect();
    window.addEventListener("keydown", handleKeyDown);
    onCleanup(() => {
      window.removeEventListener("keydown", handleKeyDown);
    });
  });

  const openModal = async (user: User) => {
    setSelectedUser(user);
    setSelectedRoles(user.roles.map((r) => r.id));

    // Get current special fields data (depends on roles)
    if (user.roles.some((r) => r.id === 3)) {
      // institution administrator
      const institutionAdministratorData = await api.get<{
        institutionAdministrator: InstitutionAdministrator;
      }>(`/institution_administrators/${user.id}`);
      setInstitutionAdministratorPositionId(
        institutionAdministratorData.institutionAdministrator.position?.id ||
          null,
      );
    }
    if (user.roles.some((r) => r.id === 4)) {
      // staff
      const staffData = await api.get<{ staff: Staff }>(`/staff/${user.id}`);
      setStaffPositionId(staffData.staff.position?.id || null);
    }
    if (user.roles.some((r) => r.id === 5)) {
      // teacher
      const teacherData = await api.get<{ teacher: Teacher }>(
        `/teachers/${user.id}`,
      );
      setTeacherClassroomId(teacherData.teacher.classroom?.id || null);
      setTeacherSubjectIds(
        teacherData.teacher.subjects?.map((s: Subject) => s.id) || [],
      );
      setTeacherStudentGroupIds(
        teacherData.teacher.studentGroups?.map((g: StudentGroup) => g.id) || [],
      );
    }
    if (user.roles.some((r) => r.id === 6)) {
      // parent
      const parentData = await api.get<{ parent: Parent }>(
        `/parents/${user.id}`,
      );
      setParentStudentIds(
        parentData.parent.students?.map((s: Student) => s.userId) || [],
      );
    }
    if (user.roles.some((r) => r.id === 7)) {
      // student
      const studentData = await api.get<{ student: Student }>(
        `/students/${user.id}`,
      );
      setStudentGroupId(studentData.student.studentGroup?.id || null);
    }

    setError("");
  };

  const closeModal = () => {
    setSelectedUser(null);
    setSelectedRoles([]);
    setTeacherClassroomId(null);
    setTeacherSubjectIds([]);
    setStudentGroupId(null);
    setStaffPositionId(null);
    setParentStudentIds([]);
    setInstitutionAdministratorPositionId(null);
    setSaving(false);
    setError("");
  };

  const toggleRole = (roleId: number) => {
    if (selectedRoles().includes(roleId)) {
      if (selectedRoles().length === 1) {
        setError("Нельзя удалить последнюю роль");
        return;
      }
      setSelectedRoles(selectedRoles().filter((id) => id !== roleId));
    } else {
      setSelectedRoles([...selectedRoles(), roleId]);
    }
    setError("");
  };

  const saveRoles = async () => {
    if (!selectedUser()) return;

    setError("");

    const formData = new URLSearchParams();
    selectedRoles().forEach((id) => formData.append("roleId", id.toString()));

    // Special fields (depends on selected roles)
    if (selectedRoles().includes(3)) {
      if (!institutionAdministratorPositionId()) {
        setError("Выберите должность администрации ОУ");
        return;
      }
      formData.append(
        "institutionAdministratorPositionId",
        String(institutionAdministratorPositionId()),
      );
    }
    if (selectedRoles().includes(4)) {
      if (!staffPositionId()) {
        setError("Выберите должность сотрудника ОУ");
        return;
      }
      formData.append("staffPositionId", String(staffPositionId()));
    }
    if (selectedRoles().includes(5)) {
      if (!teacherClassroomId()) {
        setError("Выберите аудиторию/классный кабинет");
        return;
      }
      formData.append("teacherClassroomId", String(teacherClassroomId()));
      if (teacherSubjectIds().length === 0) {
        setError("Выберите предметы");
        return;
      }
      teacherSubjectIds().forEach((subjectId) => {
        formData.append("teacherSubjectId", String(subjectId));
      });
      if (teacherStudentGroupIds.length > 0) {
        teacherStudentGroupIds.forEach((groupId) => {
          formData.append("teacherStudentGroupId", String(groupId));
        });
      }
    }
    if (selectedRoles().includes(6)) {
      parentStudentIds.forEach((studentId) => {
        if (studentId.trim()) {
          formData.append("parentStudentId", studentId);
        }
      });
    }
    if (selectedRoles().includes(7)) {
      if (!studentGroupId()) {
        setError("Выберите класс/учебную группу");
        return;
      }
      formData.append("studentGroupId", String(studentGroupId()));
    }

    try {
      setSaving(true);
      await api.put(
        `/users/${selectedUser()!.id}/roles${!hasPermission(PERMISSIONS.ROLE_ADMIN_ASSIGN) ? "/non_admin" : ""}`,
        formData,
      );
      // Update locally
      setUsers(
        users().map((u) =>
          u.id === selectedUser()!.id
            ? {
                ...u,
                roles: ROLES_TO_DISPLAY.filter((r) =>
                  selectedRoles().includes(r.id),
                ).map((r) => ({
                  id: r.id,
                  name: r.name,
                  createdAt: "", // TODO: is it OK? :)
                  updatedAt: "",
                  permissions: [],
                })),
              }
            : u,
        ),
      );
      closeModal();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка сохранения ролей");
    } finally {
      setSaving(false);
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

  const deleteUser = async (user: User) => {
    if (
      !confirm(
        `Удалить пользователя ${user.firstName} ${user.lastName}? Это действие нельзя отменить.`,
      )
    ) {
      return;
    }

    try {
      await api.delete(`/users/${user.id}`);
      setUsers(users().filter((u) => u.id !== user.id));
      if (users().length === 0 && page() > 0) {
        setPage((prev) => prev - 1);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка удаления");
    }
  };

  const [showCopied, setShowCopied] = createSignal(false);
  const copyToClipboard = async (text: string) => {
    await navigator.clipboard.writeText(text);
    setShowCopied(true);
    setTimeout(() => setShowCopied(false), 2000);
  };

  const addAdminRole = async (user: User) => {
    try {
      setSaving(true);
      const formData = new URLSearchParams();
      formData.append("roleId", "2");

      await api.post(`/users/${user.id}/roles`, formData);

      // Update locally (add admin role)
      setUsers(
        users().map((u) =>
          u.id === user.id
            ? {
                ...u,
                roles: [
                  ...u.roles,
                  {
                    id: 2,
                    name: ROLES.ADMIN,
                    createdAt: "", // TODO: is it OK? :)
                    updatedAt: "",
                    permissions: [],
                  },
                ].filter(
                  (r, index, self) =>
                    index === self.findIndex((t) => t.id === r.id), // remove duplicates
                ),
              }
            : u,
        ),
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка добавления роли админа",
      );
    } finally {
      setSaving(false);
    }
  };

  const removeAdminRole = async (user: User) => {
    try {
      setSaving(true);
      await api.delete(`/users/${user.id}/roles/2`);

      // Update locally (remove admin role)
      setUsers(
        users().map((u) =>
          u.id === user.id
            ? {
                ...u,
                roles: u.roles.filter((r) => r.id !== 2),
              }
            : u,
        ),
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка удаления роли админа",
      );
    } finally {
      setSaving(false);
    }
  };

  // TODO: mark required fields with "*"
  return (
    <div class="space-y-6 p-4">
      {/* "Copied!" notification */}
      <Show when={showCopied()}>
        <div class="fixed top-5 left-1/2 -translate-x-1/2 z-50">
          <div class="bg-gray-800 text-white px-5 py-3 rounded-xl shadow-lg text-sm font-medium">
            ✓ Скопировано!
          </div>
        </div>
      </Show>

      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">
          Управление пользователями
        </h1>
        <p class="text-gray-500 mt-1">Назначение и изменение ролей</p>
      </div>

      <Show when={error() && !selectedUser()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      <Show when={loading()}>
        <div class="flex justify-center py-16">
          <div class="text-gray-500">Загрузка...</div>
        </div>
      </Show>

      <Show when={!loading() && users().length === 0}>
        <div class="flex flex-col items-center justify-center gap-1 py-16">
          <UsersIcon class="w-15 h-15 mb-3" />
          <p class="text-gray-500">Нет пользователей</p>
        </div>
      </Show>

      <Show when={!loading() && users().length > 0}>
        <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-6 py-4 text-center text-sm font-semibold text-gray-600">
                    Пользователь
                  </th>
                  <th class="px-6 py-4 text-center text-sm font-semibold text-gray-600">
                    Email
                  </th>
                  <th class="px-6 py-4 text-center text-sm font-semibold text-gray-600">
                    ID
                  </th>
                  <th class="px-6 py-4 text-center text-sm font-semibold text-gray-600">
                    Роли
                  </th>
                  <th class="px-6 py-4 text-center text-sm font-semibold text-gray-600">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <For
                  each={users().filter((user) =>
                    user.roles.every((role) => role.id !== 1),
                  )}
                >
                  {(user) => (
                    <tr class="hover:bg-gray-50 transition">
                      <td class="px-6 py-4">
                        <div class="flex items-center gap-3">
                          <A
                            href={`/users/${user.id}`}
                            class="w-8 h-8 bg-gray-100 rounded-full hover:bg-gray-200 transition aspect-square"
                          >
                            <img
                              class="w-8 h-8 rounded-full object-cover border-2 border-blue-100 hover:brightness-95 transition"
                              src={`/storage/storage/avatars/${user.hasAvatar ? user.id : "default"}.jpeg`}
                              alt="Фото профиля"
                            />
                          </A>
                          <span
                            class={`${auth.user()?.id === user.id ? "font-semibold" : ""}`}
                          >
                            {user.lastName} {user.firstName}{" "}
                            {user?.middleName || ""}
                          </span>
                        </div>
                      </td>
                      <td
                        class="px-6 py-4 text-sm text-gray-500 cursor-copy"
                        onClick={() => {
                          copyToClipboard(user.email);
                        }}
                      >
                        {user.email}
                      </td>
                      <td
                        class="px-6 py-4 text-sm text-gray-400 font-mono cursor-copy"
                        onClick={() => {
                          copyToClipboard(user.id);
                        }}
                      >
                        {user.id.slice(0, 8)}...
                      </td>
                      <td class="px-6 py-4">
                        <div class="flex flex-wrap gap-1.5">
                          <For each={user.roles}>
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
                      </td>
                      <td class="px-6 py-4 space-x-3 flex align-center justify-center">
                        <Show
                          when={hasPermission(PERMISSIONS.ROLE_USER_ASSIGN)}
                        >
                          <button
                            onClick={() => openModal(user)}
                            class="text-blue-600 hover:text-blue-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
                          >
                            Изменить роли
                          </button>
                        </Show>
                        <Show
                          when={hasPermission(PERMISSIONS.ROLE_ADMIN_ASSIGN)}
                        >
                          <Show
                            when={user.roles.every(
                              (r) => r.name !== ROLES.ADMIN,
                            )}
                          >
                            <button
                              onClick={() => addAdminRole(user)}
                              class="text-blue-600 hover:text-blue-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
                            >
                              Сделать админом
                            </button>
                          </Show>
                          <Show
                            when={
                              user.roles.some((r) => r.name === ROLES.ADMIN) &&
                              user.roles.length > 1
                            }
                          >
                            <button
                              onClick={() => removeAdminRole(user)}
                              class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
                            >
                              Снять права админа
                            </button>
                          </Show>
                        </Show>
                        <Show
                          when={
                            (user.roles.every(
                              (r) =>
                                r.name !== ROLES.ADMIN &&
                                r.name !== ROLES.SUPERADMIN,
                            ) &&
                              hasPermission(
                                PERMISSIONS.USER_DELETE_ANY_USER,
                              )) ||
                            (user.roles.every((r) => r.name === ROLES.ADMIN) &&
                              hasPermission(PERMISSIONS.USER_DELETE_ANY_ADMIN))
                          }
                        >
                          <button
                            onClick={() => deleteUser(user)}
                            class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed
"
                          >
                            Удалить
                          </button>
                        </Show>
                      </td>
                    </tr>
                  )}
                </For>
              </tbody>
            </table>
          </div>
        </div>
      </Show>

      {/* Modal */}
      <Show when={selectedUser()}>
        <div
          class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
          onClick={closeModal}
        >
          <div
            class="bg-white rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-hidden"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4">
              <div class="flex items-center gap-3">
                <img
                  src={`/storage/storage/avatars/${selectedUser()?.hasAvatar ? selectedUser()?.id : "default"}.jpeg`}
                  alt="Аватар"
                  class="w-10 h-10 rounded-full object-cover"
                />
                <div>
                  <h2 class="text-xl font-bold text-gray-800">
                    {selectedUser()?.lastName} {selectedUser()?.firstName}
                  </h2>
                  <p class="text-sm text-gray-500">{selectedUser()?.email}</p>
                </div>
              </div>
            </div>

            {/* Body */}
            <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)] space-y-5">
              <Show when={error()}>
                <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl text-sm">
                  {error()}
                </div>
              </Show>

              {/* Roles selection */}
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Выберите роли:
                </label>
                <div class="space-y-2">
                  <For
                    each={ROLES_TO_DISPLAY.filter((role) => {
                      if (hasRole(ROLES.ADMIN)) {
                        return (
                          role.name !== ROLES.SUPERADMIN &&
                          role.name !== ROLES.ADMIN
                        ); // TODO: make in the whole frontend code like here
                      }
                      return false;
                    })}
                  >
                    {(role) => (
                      <label class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition cursor-pointer">
                        <input
                          type="checkbox"
                          checked={selectedRoles().includes(role.id)}
                          onChange={() => {
                            setSaving(false);
                            setError("");
                            toggleRole(role.id);
                          }}
                          disabled={saving()}
                          class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                        />
                        <span class="text-gray-700">{role.displayName}</span>
                      </label>
                    )}
                  </For>
                </div>
              </div>

              {/* Institution Administrator Position */}
              <Show when={selectedRoles().includes(3)}>
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

              {/* Staff Position */}
              <Show when={selectedRoles().includes(4)}>
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

              {/* Teacher fields */}
              <Show when={selectedRoles().includes(5)}>
                <div class="space-y-4 border-t border-gray-100 pt-4">
                  <h3 class="font-medium text-gray-800">Данные учителя</h3>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Классный кабинет *
                    </label>
                    <select
                      disabled={saving()}
                      value={teacherClassroomId() || ""}
                      onChange={(e) => {
                        setSaving(false);
                        setError("");
                        setTeacherClassroomId(Number(e.currentTarget.value));
                      }}
                      class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
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
                    <div class="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto border border-gray-200 rounded-xl p-2">
                      <For each={subjects()}>
                        {(subject) => (
                          <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded-lg cursor-pointer transition">
                            <input
                              disabled={saving()}
                              type="checkbox"
                              checked={teacherSubjectIds().includes(subject.id)}
                              onChange={() => {
                                setSaving(false);
                                setError("");
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
                            class="max-md:aspect-square flex items-center justify-center px-2 md:px-4 bg-red-700 text-white rounded-xl hover:bg-red-800 transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                          >
                            <span class="hidden md:flex">Удалить</span>
                            <Trash class="flex md:hidden" />
                          </button>
                        </div>
                      )}
                    </Index>

                    <button
                      type="button"
                      onClick={addTeacherStudentGroupId}
                      class="w-full py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer flex flex-row flex-nowrap items-center justify-center gap-2"
                    >
                      <Plus /> Добавить группу
                    </button>
                  </div>
                </div>
              </Show>

              {/* Parent fields */}
              <Show when={selectedRoles().includes(6)}>
                <div class="space-y-3 border-t border-gray-100 pt-4">
                  <h3 class="font-medium text-gray-800">Привязка учеников</h3>

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
                              class="max-md:aspect-square flex items-center justify-center px-2 md:px-4 bg-red-700 text-white rounded-xl hover:bg-red-800 transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                            >
                              <span class="hidden md:flex">Удалить</span>
                              <Trash class="flex md:hidden" />
                        </button>
                      </div>
                    )}
                  </Index>

                  <button
                    disabled={saving()}
                    type="button"
                    onClick={addStudentId}
                        class="w-full py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer flex flex-row flex-nowrap items-center justify-center gap-2"
                      >
                        <Plus /> Добавить ученика
                  </button>
                </div>
              </Show>

              {/* Student fields */}
              <Show when={selectedRoles().includes(7)}>
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

            {/* Footer */}
            <div class="sticky bottom-0 bg-white border-t border-gray-200 px-6 py-4 flex justify-end gap-3">
              <button
                onClick={closeModal}
                class="w-40 h-10 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
              >
                Отмена
              </button>
              <button
                onClick={saveRoles}
                disabled={saving()}
                class="w-40 h-10 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
              >
                Сохранить
              </button>
            </div>
          </div>
        </div>
      </Show>
      <Pagination
        page={page()}
        hasMore={hasMore()}
        onPrev={() => setPage((prev) => prev - 1)}
        onNext={() => setPage((prev) => prev + 1)}
      />
    </div>
  );
};

export default Users;
