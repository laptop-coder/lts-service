import { createSignal, Show, For, Index, onMount } from "solid-js";
import { createStore } from "solid-js/store";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions, ROLES_TO_DISPLAY } from "../../lib/permissions";
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
} from "../../lib/types";

const Users = () => {
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
  const [studentGroupId, setStudentGroupId] = createSignal<number | null>(null);
  const [staffPositionId, setStaffPositionId] = createSignal<number | null>(
    null,
  );
  const [
    institutionAdministratorPositionId,
    setInstitutionAdministratorPositionId,
  ] = createSignal<number | null>(null);
  const [parentStudentIds, setParentStudentIds] = createStore<string[]>([]);

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

  const { hasPermission } = usePermissions();

  const loadUsers = async () => {
    try {
      const data = await api.get<{ users: User[] }>("/users");
      setUsers(data.users);
    } catch (err) {
      setError(
        "Ошибка загрузки пользователей", // TODO
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
      api.get<{ rooms: Room[] }>("/rooms"),
      api.get<{ subjects: Subject[] }>("/subjects"),
      api.get<{ studentGroups: StudentGroup[] }>("/student_groups"),
      api.get<{ staffPositions: StaffPosition[] }>("/staff/positions"),
      api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions"),
    ]);
    setRooms(roomsData.rooms);
    setSubjects(subjectsData.subjects);
    setStudentGroups(groupsData.studentGroups);
    setStaffPositions(staffPositionData.staffPositions);
    setInstitutionAdministratorPositions(
      institutionAdministratorPositionData.institutionAdministratorPositions,
    );
  };

  onMount(() => {
    loadUsers();
    loadDataForSelect();
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
    setSaving(false)
    setError("")
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
      console.log(formData.toString())
      await api.put(`/users/${selectedUser()!.id}/roles`, formData);
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
      setError("Ошибка сохранения ролей");
    } finally {
      setSaving(false);
    }
  };

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

  return (
    <div class="space-y-6">
      <h1 class="text-2xl font-bold">Управление ролями пользователей</h1>

      <Show when={error() && !selectedUser()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={!loading() && users().length === 0}>
        <div class="text-center text-gray-500 py-8">Нет пользователей</div>
      </Show>

      <Show when={!loading() && users().length > 0}>
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                  ФИО
                </th>
                <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                  Email
                </th>
                <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                  Роли
                </th>
                <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <For
                each={users().filter((user) =>
                  user.roles.every((role) => role.id !== 1),
                )}
              >
                {(user) => (
                  <tr class="hover:bg-gray-50">
                    <td class="px-4 py-3">
                      {user.lastName} {user.firstName} {user?.middleName}
                    </td>
                    <td class="px-4 py-3 text-sm text-gray-500">
                      {user.email}
                    </td>
                    <td class="px-4 py-3">
                      <div class="flex flex-wrap gap-1">
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
                    <td class="px-4 py-3 text-right">
                      <Show when={hasPermission(PERMISSIONS.ROLE_ASSIGN)}>
                        <button
                          onClick={() => openModal(user)}
                          class="text-blue-600 hover:text-blue-800"
                        >
                          Изменить роли
                        </button>
                      </Show>
                    </td>
                  </tr>
                )}
              </For>
            </tbody>
          </table>
        </div>
      </Show>

      {/* Modal */}
      <Show when={selectedUser()}>
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
            <div class="p-6">
              <h2 class="text-xl font-bold mb-4">
                Роли пользователя: {selectedUser()?.firstName}{" "}
                {selectedUser()?.lastName}
              </h2>

              <Show when={error()}>
                <div class="bg-red-100 text-red-700 p-2 rounded mb-4 text-sm">
                  {error()}
                </div>
              </Show>

              <div class="space-y-2 mb-6">
                <div class="text-sm font-medium text-gray-700 mb-2">
                  Выберите роли:
                </div>
                <For
                  each={ROLES_TO_DISPLAY.filter(
                    (role) => ![1, 2].includes(role.id),
                  )}
                >
                  {(role) => (
                    <label class="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={selectedRoles().includes(role.id)}
                        onChange={() => {
                          setSaving(false);
                          setError("");
                          toggleRole(role.id);
                        }}
                        disabled={saving()}
                      />
                      <span>{role.displayName}</span>
                    </label>
                  )}
                </For>
              </div>
              <Show when={selectedRoles().includes(3)}>
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
                >
                  <option value="">Выберете должность администрации ОУ*</option>
                  <For each={institutionAdministratorPositions()}>
                    {(position) => (
                      <option value={position.id}>{position.name}</option>
                    )}
                  </For>
                </select>
              </Show>
              <Show when={selectedRoles().includes(4)}>
                <select
                  disabled={saving()}
                  value={staffPositionId() || ""}
                  onChange={(e) => {
                    setSaving(false);
                    setError("");
                    setStaffPositionId(Number(e.currentTarget.value));
                  }}
                >
                  <option value="">Выберете должность сотрудника ОУ*</option>
                  <For each={staffPositions()}>
                    {(position) => (
                      <option value={position.id}>{position.name}</option>
                    )}
                  </For>
                </select>
              </Show>
              <Show when={selectedRoles().includes(5)}>
                <select
                  disabled={saving()}
                  value={teacherClassroomId() || ""}
                  onChange={(e) => {
                    setSaving(false);
                    setError("");
                    setTeacherClassroomId(Number(e.currentTarget.value));
                  }}
                >
                  <option value="">Выберете аудиторию/классный кабинет*</option>
                  <For each={rooms()}>
                    {(room) => <option value={room.id}>{room.name}</option>}
                  </For>
                </select>
                Выберете предметы
                <For each={subjects()}>
                  {(subject) => (
                    <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded cursor-pointer transition">
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
                      <span class="text-gray-700">{subject.name}</span>
                    </label>
                  )}
                </For>
              </Show>
              <Show when={selectedRoles().includes(6)}>
                <div class="space-y-2">
                  <label class="block text-sm font-medium text-gray-700">
                    ID учеников (для привязки к родителю)*
                  </label>
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
                          class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                        <button
                  disabled={saving()}
                          type="button"
                          onClick={() => removeStudentId(index)}
                          class="px-3 py-2 bg-red-500 text-white rounded-md hover:bg-red-600"
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
                    class="w-full px-3 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition"
                  >
                    + Добавить ученика
                  </button>
                </div>
              </Show>
              <Show when={selectedRoles().includes(7)}>
                <select
                  disabled={saving()}
                  value={studentGroupId() || ""}
                  onChange={(e) =>
                    setStudentGroupId(Number(e.currentTarget.value))
                  }
                >
                  <option value="">Выберете класс/учебную группу*</option>
                  <For each={studentGroups()}>
                    {(group) => <option value={group.id}>{group.name}</option>}
                  </For>
                </select>
              </Show>

              <div class="flex justify-end gap-2">
                <button
                  onClick={closeModal}
                  class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                >
                  Отмена
                </button>
                <button
                  onClick={saveRoles}
                  disabled={saving()}
                  class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  {saving() ? "Сохранение..." : "Сохранить"}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default Users;
