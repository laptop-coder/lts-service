import axiosInstance from './axiosInstance';

const thingMarkAsFound = async (props: { thing: { id: string } }) => {
  return axiosInstance
    .post(
      `/thing/mark_as_found`, // TODO: move to 'utils/consts'
      { thingId: props.thing.id },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then(() => {
      // window.location.href = HOME__ROUTE;
    })
    .catch((error) =>
      // console.log(error) // TODO: think about it, how to make right
      console.log('error'),
    );
};

export default thingMarkAsFound;
