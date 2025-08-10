export const GET = (path: string) => {
  return new Promise((resolve, reject) => {
    fetch(`https://172.16.1.2/${path}`)
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
