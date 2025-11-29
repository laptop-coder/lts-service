const checkPhotoAvailability = (props: {
  pathToPhoto: string;
  success: () => boolean;
  // failure: () => boolean; // you can use it if you need
}): void => {
  /**
   * This function changes the value of the signal if the photo is available.
   *
   * @param pathToPhoto<string> - Path to photo on the server.
   * @param success<() => boolean> - Signal setter (e. g., `() => set(true)`).
   * Performed if the photo is available.
   * @returns <undefined> - The function doesn't returns anything.
   */
  var photo = new Image();
  photo.onload = props.success;
  // photo.onerror = props.failure;
  photo.src = props.pathToPhoto;
};

export default checkPhotoAvailability;
