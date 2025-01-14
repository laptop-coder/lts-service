export const POST = async (path: string, data) => {
  const response = await fetch(`http://localhost:8000/${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json; charset=utf-8",
    },
    body: JSON.stringify(data),
  });
  return response.json();
};
