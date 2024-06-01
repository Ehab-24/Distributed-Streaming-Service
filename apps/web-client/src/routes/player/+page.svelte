<script lang="ts">
  import { onMount } from "svelte";
  import shaka from "shaka-player";
  import { page } from "$app/stores";

  const isInteger = (str: string | null): boolean =>
    str !== null && /^\d+$/.test(str);

  let manifestUri: string;
  let player: any;

  function getServerPort(server_id: string): number {
    switch (server_id) {
      case "4":
        return 5000;
      case "5":
        return 5001;
      case "6":
        return 5002;
      default:
        throw new Error("Invalid server_id");
    }
  }

  let error: string | null = null;

  onMount(() => {
    const video_id = $page.url.searchParams.get("video_id") ?? "";
    if (!isInteger(video_id)) {
      error = "Query params: Invalid video_id";
      return;
    }
    let chunk_id = $page.url.searchParams.get("chunk_id") ?? "";
    if (!isInteger(chunk_id)) {
      error = "Query params: Invalid chunk_id";
      return;
    }
    const server_id = $page.url.searchParams.get("server_id") ?? "";
    if (!isInteger(server_id)) {
      error = "Query params: Invalid server_id";
      return;
    }

    function getManifestUri(chunk_id: string): string {
      return `http://127.0.0.1:${getServerPort(server_id)}/data/${server_id}/media/processed/${video_id}_${chunk_id}/720p.mpd`;
    }
    manifestUri = getManifestUri(chunk_id);

    // const manifestUri =
    //   `http://127.0.0.1:${getServerPort(server_id)}/data/${server_id}/media/processed/merged.mpd`;
    function initApp() {
      shaka.polyfill.installAll();
      if (shaka.Player.isBrowserSupported()) {
        console.log("Initializing player");
        initPlayer();
      } else {
        console.error("Browser not supported!");
      }
    }

    async function loadPlayer(uri: string) {
      try {
        await player.load(uri);
        console.log("The video has now been loaded!");
      } catch (e) {
        onError(e);
      }
    }

    async function initPlayer() {
      const video = document.getElementById("video");
      video?.addEventListener("ended", onEnded)
      player = new shaka.Player();
      await player.attach(video);
      window.player = player;
      loadPlayer(manifestUri);
    }

    function onError(error: any) {
      console.error("Error code", error.code, "object", error);
    }
    function onEnded(ev: Event) {
      console.log("Video ended");
      chunk_id = (parseInt(chunk_id) + 1).toString();
      manifestUri = getManifestUri(chunk_id);
      loadPlayer(manifestUri);
    }

    initApp();
  });
</script>

{#if error}
  <p>{error}</p>
{:else}
  <video
    id="video"
    width="100%"
    poster="//shaka-player-demo.appspot.com/assets/poster.jpg"
    controls
    autoplay
  ></video>
{/if}
