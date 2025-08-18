import { createSignal } from 'solid-js';

const checkPhotoAvailability = (props: { pathToPhoto: string }): boolean => {
  /**
   * This function checks if the photo is available and returns the result
   *
   * @param pathToPhoto<string> - Path to photo on the server.
   * @returns <boolean> - Returns true if the photo on the server is available
   * and false if not.
   */
  const [thingPhotoIsAvailable, setThingPhotoIsAvailable] = createSignal(false);
  const photo = new Image();
  photo.onload = () => setThingPhotoIsAvailable(true);
  photo.src = props.pathToPhoto;
  return thingPhotoIsAvailable();
};

export default checkPhotoAvailability;
