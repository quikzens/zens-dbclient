<script>
  import { Router, Route } from "svelte-routing";
  import TablePage from "./page/TablePage.svelte";
  import RecordPage from "./page/RecordPage.svelte";
  import ConnectionPage from "./page/ConnectionPage.svelte";

  export let url = "/table";
</script>

<Router {url}>
  <div>
    <Route path="/record">
      <RecordPage />
    </Route>
    <Route path="/connection">
      <ConnectionPage />
    </Route>
    <Route path="/table">
      <TablePage />
    </Route>
  </div>
</Router>
