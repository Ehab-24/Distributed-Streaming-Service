<script lang="ts">
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import { PUBLIC_SERVER_URL } from "$env/static/public";

  let videos: any[] = [];
  let searchQuery = writable("");

  const serverUrl = (path: string = "/") => `${PUBLIC_SERVER_URL}${path}`;

  onMount(async () => {
    const response = await fetch(serverUrl("/all/"));
    videos = await response.json();
  });

  $: filteredVideos = videos.filter((video) =>
    video.title.toLowerCase().includes($searchQuery.toLowerCase()),
  );
</script>

<div class="container mx-auto px-4 py-6">
  <div class="text-center mb-8">
    <h1 class="text-4xl font-bold mb-2">Video Gallery</h1>
    <p class="text-lg text-gray-600">
      Discover and explore our collection of videos
    </p>
  </div>

  <div class="mb-6">
    <input
      type="text"
      placeholder="Search for videos..."
      class="w-full p-3 rounded border border-gray-300 focus:outline-none focus:border-blue-500"
      bind:value={$searchQuery}
    />
  </div>

  <div
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
  >
    {#each filteredVideos as video}
      <a
        href={`/player?video_id=${video.id}&chunk_id=${video.first_chunk.id}&server_id=${video.first_chunk.server_id}`}
        class="max-w-xs"
      >
        <div
          class="w-full aspect-[16/10] rounded-xl overflow-hidden"
        >
          <img
            class="w-full h-full object-fit bg-gray-300"
            src="https://loremflickr.com/320/240?random=1"
            alt={video.title}
          />
        </div>
        <div class="font-bold text-md pt-4 pb-2 break-words">{video.title}</div>
        <div class="flex gap-4 text-sm font-light">
          <span
            class="inline-block text-gray-600 rounded-full text-sm font-semibold"
            >{video.duration}s</span
          >
          <span
            class="inline-block text-gray-600 rounded-full text-sm font-semibold"
            >{video.replication_factor} replicas</span
          >
        </div>
      </a>
    {/each}
  </div>
</div>

<style>
  .container {
    max-width: 1200px;
  }
</style>

