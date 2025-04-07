export const POST = (path: string, data) => {
  return new Promise((resolve, reject) => {
    fetch(`http://localhost:8080/${path}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json; charset=utf-8",
      },
      body: JSON.stringify(data),
    })
      .then((response) => {
        if (!response.ok) {
          reject(new Error(`Error! Status: ${response.status}`));
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};
