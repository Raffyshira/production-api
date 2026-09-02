export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/v1";

export const fetcher = (at: string) => (url: string) =>
  fetch(API_URL + url, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${at}`,
    },
  }).then((r) => r.json());
