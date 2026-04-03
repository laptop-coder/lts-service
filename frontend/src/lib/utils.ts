export const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString("ru-RU");
};
