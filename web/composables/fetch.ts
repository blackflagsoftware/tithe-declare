import { ofetch } from 'ofetch';

export const apiFetch = () => {
  const config = useRuntimeConfig();
  const fetch = ofetch.create({
    // Set the base URL from the runtime config
    baseURL: config.public.apiURL,
    // Handle API errors globally
    onResponseError({ request, response, options }) {
      console.error('[fetch error]', response.status, response.statusText, response._data);
	  // TODO: toast or some kind of user notification
    }
  });

  return fetch
}