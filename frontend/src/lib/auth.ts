import { createSignal } from "solid-js";
import { api } from "./api";
import type { User } from "./types";

const [user, setUser] = createSignal<User | null>(null);
const [isLoading, setIsLoading] = createSignal(true);

export function useAuth() {
  // Function to register and load user
  const register = async (formData: FormData) => {
    const data = await api.post<{ user: User }>("/users", formData);
    setUser(data.user);
    return data.user;
  };

  // Function to log in and load user
  const login = async (email: string, password: string) => {
    const formData = new URLSearchParams();
    formData.append("email", email);
    formData.append("password", password);
    const data = await api.post<{ user: User }>("/auth/login", formData);
    setUser(data.user);
    return data.user;
  };

  // Function to log out and clear user
  const logout = async () => {
    await api.post("/auth/logout");
    setUser(null);
  };

  // Function to load user if logged in
  const checkAuth = async () => {
    try {
      const data = await api.get<{ user: User }>("/users/me");
      setUser(data.user);
    } catch {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    user,
    isLoading,
    login,
    register,
    logout,
    checkAuth,
    isAuthenticated: () => user() !== null, //TODO: is it necessary?
  };
}
