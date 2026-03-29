import { createSignal, onMount, Show, For, Index } from "solid-js";
import { createStore } from "solid-js/store";
import { useAuth } from "../lib/auth";
import { api } from "../lib/api";
import { useNavigate, useSearchParams } from "@solidjs/router";
import type {
  Role,
  Room,
  Subject,
  StudentGroup,
  StaffPosition,
  InstitutionAdministratorPosition,
} from "../lib/types";

// TODO: add avatar support
const Register = () => {
  // Data about new user
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [firstName, setFirstName] = createSignal("");
  const [middleName, setMiddleName] = createSignal("");
  const [lastName, setLastName] = createSignal("");
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

  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const [searchParams] = useSearchParams();
  const auth = useAuth();
  const navigate = useNavigate();

  const inviteToken = searchParams.inviteToken;
  if (typeof inviteToken !== "string") {
    return <div>Для регистрации нужен инвайт-токен</div>;
  }
  // Roles
  const [roleIds, setRoleIds] = createSignal<number[]>([]);
  const [roleNames, setRoleNames] = createSignal<string[]>([]);
  // Rooms
  const [rooms, setRooms] = createSignal<Room[]>([]);
  // Subjects
  const [subjects, setSubjects] = createSignal<Subject[]>([]);
  // Student groups
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  // Staff positions
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  // Institution administrators positions
  const [
    institutionAdministratorPositions,
    setInstitutionAdministratorPositions,
  ] = createSignal<InstitutionAdministratorPosition[]>([]);

  onMount(async () => {
    try {
      // Get roles from the invite token
      const rolesData = await api.get<{ roles: Role[] }>(
        `/tokens/invite/${inviteToken}`,
      );
      setRoleIds(rolesData.roles.map((role) => role.id));
      setRoleNames(rolesData.roles.map((role) => role.name));
      // Get all rooms
      const roomsData = await api.get<{ rooms: Room[] }>("/rooms");
      setRooms(roomsData.rooms);
      // Get all subjects
      const subjectsData = await api.get<{ subjects: Subject[] }>("/subjects");
      setSubjects(subjectsData.subjects);
      // Get all student groups
      const studentGroupsData = await api.get<{
        studentGroups: StudentGroup[];
      }>("/student_groups");
      setStudentGroups(studentGroupsData.studentGroups);
      // Get all staff positions
      const staffPositionsData = await api.get<{
        staffPositions: StaffPosition[];
      }>("/staff/positions");
      setStaffPositions(staffPositionsData.staffPositions);
      // Get all institution administrator positions
      const institutionAdministratorPositionsData = await api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions");
      setInstitutionAdministratorPositions(
        institutionAdministratorPositionsData.institutionAdministratorPositions,
      );
    } catch (err) {
      // TODO
    } finally {
      setLoading(false);
    }
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    const formData = new FormData();
    formData.append("email", email());
    formData.append("password", password());
    formData.append("firstName", firstName());
    if (middleName()?.trim()) formData.append("middleName", middleName());
    formData.append("lastName", lastName());
    formData.append("inviteToken", inviteToken);
    if (roleIds().includes(3)) {
      if (!institutionAdministratorPositionId()) {
        setError("Выберите должность администрации ОУ");
        return;
      }
      formData.append(
        "institutionAdministratorPositionId",
        String(institutionAdministratorPositionId()),
      );
    }
    if (roleIds().includes(4)) {
      if (!staffPositionId()) {
        setError("Выберите должность сотрудника ОУ");
        return;
      }
      formData.append("staffPositionId", String(staffPositionId()));
    }
    if (roleIds().includes(5)) {
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
    if (roleIds().includes(6)) {
      parentStudentIds.forEach((studentId) => {
        if (studentId.trim()) {
          formData.append("parentStudentId", studentId);
        }
      });
    }
    if (roleIds().includes(7)) {
      if (!studentGroupId()) {
        setError("Выберите класс/учебную группу");
        return;
      }
      formData.append("studentGroupId", String(studentGroupId()));
    }

    try {
      await auth.register(formData);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка регистрации",
      );
    } finally {
      setLoading(false);
      navigate("/");
    }
  };

  const addStudentId = () => {
    setParentStudentIds([...parentStudentIds, ""]);
  };

  const updateStudentId = (index: number, value: string) => {
    setParentStudentIds(index, value);
  };

  const removeStudentId = (index: number) => {
    setParentStudentIds(parentStudentIds.filter((_, i) => i !== index));
  };

  return (
    <>
      {loading() ? (
        <div>Загрузка...</div>
      ) : (
        <form
          onSubmit={handleSubmit}
          class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6 space-y-4"
        >
          Регистрация пользователя со следующими ролями: {roleNames()}
          <input
            type="email"
            value={email()}
            onInput={(e) => setEmail(e.currentTarget.value)}
            placeholder="Email*"
            required
          />
          <input
            type="password"
            value={password()}
            onInput={(e) => setPassword(e.currentTarget.value)}
            placeholder="Пароль*"
            required
          />
          <input
            type="text"
            value={firstName()}
            onInput={(e) => setFirstName(e.currentTarget.value)}
            placeholder="Имя*"
            required
          />
          <input
            type="text"
            value={middleName()}
            onInput={(e) => setMiddleName(e.currentTarget.value)}
            placeholder="Отчество"
          />
          <input
            type="text"
            value={lastName()}
            onInput={(e) => setLastName(e.currentTarget.value)}
            placeholder="Фамилия*"
            required
          />
          <Show when={roleIds().includes(3)}>
            <select
              value={institutionAdministratorPositionId() || ""}
              onChange={(e) =>
                setInstitutionAdministratorPositionId(
                  Number(e.currentTarget.value),
                )
              }
            >
              <option value="">Выберете должность администрации ОУ*</option>
              <For each={institutionAdministratorPositions()}>
                {(position) => (
                  <option value={position.id}>{position.name}</option>
                )}
              </For>
            </select>
          </Show>
          <Show when={roleIds().includes(4)}>
            <select
              value={staffPositionId() || ""}
              onChange={(e) =>
                setStaffPositionId(Number(e.currentTarget.value))
              }
            >
              <option value="">Выберете должность сотрудника ОУ*</option>
              <For each={staffPositions()}>
                {(position) => (
                  <option value={position.id}>{position.name}</option>
                )}
              </For>
            </select>
          </Show>
          <Show when={roleIds().includes(5)}>
            <select
              value={teacherClassroomId() || ""}
              onChange={(e) =>
                setTeacherClassroomId(Number(e.currentTarget.value))
              }
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
                    type="checkbox"
                    checked={teacherSubjectIds().includes(subject.id)}
                    onChange={() => {
                      if (teacherSubjectIds().includes(subject.id)) {
                        setTeacherSubjectIds((prev) =>
                          prev.filter((id) => id !== subject.id),
                        );
                      } else {
                        setTeacherSubjectIds((prev) => [...prev, subject.id]);
                      }
                    }}
                    class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                  />
                  <span class="text-gray-700">{subject.name}</span>
                </label>
              )}
            </For>
          </Show>
          <Show when={roleIds().includes(6)}>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700">
                ID учеников (для привязки к родителю)*
              </label>
              <Index each={parentStudentIds}>
                {(studentId, index) => (
                  <div class="flex gap-2">
                    <input
                      type="text"
                      value={studentId()}
                      onInput={(e) => updateStudentId(index, e.target.value)}
                      placeholder={`ID ученика ${index + 1}`}
                      class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <button
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
                type="button"
                onClick={addStudentId}
                class="w-full px-3 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition"
              >
                + Добавить ученика
              </button>
            </div>
          </Show>
          <Show when={roleIds().includes(7)}>
            <select
              value={studentGroupId() || ""}
              onChange={(e) => setStudentGroupId(Number(e.currentTarget.value))}
            >
              <option value="">Выберете класс/учебную группу*</option>
              <For each={studentGroups()}>
                {(group) => <option value={group.id}>{group.name}</option>}
              </For>
            </select>
          </Show>
          {}
          {error() && <div class="error">{error()}</div>}
          <button type="submit" disabled={loading()}>
            {loading() ? "Регистрация..." : "Зарегистрироваться"}
          </button>
        </form>
      )}
    </>
  );
};

export default Register;
