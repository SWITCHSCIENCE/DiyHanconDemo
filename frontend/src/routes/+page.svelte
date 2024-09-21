<script lang="ts">
  import { onMount } from "svelte";
  import { ProgressRadial } from "@skeletonlabs/skeleton";
  import * as app from "$lib/wailsjs/go/main/App.js";
  let greet: string;
  let steer = 0;
  let force = 0;
  onMount(async () => {
    greet = await app.Greet("No Name");
    runtime.EventsOn("steer", (s) => (steer = s));
    runtime.EventsOn("force", (s) => (force = s));
  });
</script>

<div class="container h-full mx-auto flex justify-center items-center">
  <div style="transform: rotate({90 * force - 18}deg); position:relative">
    <ProgressRadial
      value={10}
      stroke={80}
      meter="stroke-secondary-500"
      track="stroke-secondary-500/30"
      strokeLinecap="round"
    />
  </div>
  <div style="transform: rotate({540 * steer - 18}deg); position: absolute">
    <ProgressRadial
      value={10}
      stroke={40}
      meter="stroke-primary-500"
      track="stroke-primary-500/30"
      strokeLinecap="round"
    />
  </div>
</div>
