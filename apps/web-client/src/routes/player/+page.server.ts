import type { PageServerLoad } from './$types';
import { PUBLIC_SERVER_URL } from '$env/static/public';

type Server = {
  id: number;
  ip: string;
  port: number;
}

export const load: PageServerLoad = async () => {
  const apiUrl = `${PUBLIC_SERVER_URL}/servers/active/`;
  const response = await fetch(apiUrl);
  if (!response.ok) {
    throw new Error(`Failed to fetch posts: ${response.statusText}`);
  }
  const data: { servers: Server[] } = await response.json();
  console.log("SERVERS:", data.servers)
  return {
    servers: data.servers
  };
};

