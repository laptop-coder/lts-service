const formatDate = (isoDateString: string) => {
  const isoDate = new Date(isoDateString);

  const day = isoDate.getDate();
  const month = isoDate.getMonth();
  const year = isoDate.getFullYear();
  const hours = isoDate.getHours();
  const minutes = isoDate.getMinutes();

  const months = [
    'января',
    'февраля',
    'марта',
    'апреля',
    'мая',
    'июня',
    'июля',
    'августа',
    'сентября',
    'октября',
    'ноября',
    'декабря',
  ];

  const formattedDate = `${day} ${months[month]} ${year}, ${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;

  return formattedDate;
};

export default formatDate;
