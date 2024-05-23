<script lang="ts">
    import { onMount } from "svelte";
    import shaka from "shaka-player";

    onMount(() => {
        const manifestUri =
            "http://127.0.0.1:5000/media/processed/1_720p.mpd";
        function initApp() {
            shaka.polyfill.installAll();
            if (shaka.Player.isBrowserSupported()) {
                console.log("Initializing player");
                initPlayer();
            } else {
                console.error("Browser not supported!");
            }
        }

        async function initPlayer() {
            const video = document.getElementById("video");
            const player = new shaka.Player();
            await player.attach(video);
            window.player = player;
            player.addEventListener("error", onErrorEvent);

            try {
                await player.load(manifestUri);
                console.log("The video has now been loaded!");
            } catch (e) {
                onError(e);
            }
        }

        function onErrorEvent(event) {
            onError(event.detail);
        }
        function onError(error) {
            console.error("Error code", error.code, "object", error);
        }

        initApp();
    });
</script>

<video
    id="video"
    width="100%"
    poster="//shaka-player-demo.appspot.com/assets/poster.jpg"
    controls
    autoplay
></video>
