const getCookie = (key: string): string | undefined => {
  const matches = document.cookie.match(
    new RegExp(
      '(?:^|; )' +
        key.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') +
        '=([^;]*)',
    ),
  );
  return matches ? decodeURIComponent(matches[1]) : undefined;
};

export default getCookie;
