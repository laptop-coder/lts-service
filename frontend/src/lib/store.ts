import { createSignal } from "solid-js";
import { conversationApi } from "../lib/api";

export const [unreadMessagesCount, setUnreadMessagesCount] = createSignal(0);

export const refreshUnreadMessagesCount = async () => {
  try {
    const data = await conversationApi.getTotalUnreadCount();
    setUnreadMessagesCount(data.unreadCount);
  } catch (err) {}
};
