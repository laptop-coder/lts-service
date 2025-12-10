import axiosInstance from './axiosInstance';

const deleteThingPhoto = async (props: { thing: { id: string } }) =>
  axiosInstance
    .post(
      `/thing/delete_photo`,
      { thingId: props.thing.id },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then((response) => {
      return response.data;
    })
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default deleteThingPhoto;
