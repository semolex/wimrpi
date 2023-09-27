<script type="ts">
    import {GeneratePortTreeForIP} from "../wailsjs/go/main/App";
    import {main} from "../wailsjs/go/models";
    import {Column, Row, Tile, TextInput, Button, TreeView, StructuredListSkeleton} from "carbon-components-svelte";
    import ArrowRight from "carbon-icons-svelte/lib/ArrowRight.svelte";
    import Add from "carbon-icons-svelte/lib/Add.svelte";

    let ips: string[] = [""];
    let portTree: main.PortTree[] = [];
    let loading = false;

    function generatePortTree(ip: string | number) {
        ip = ip.toString();
        loading = true;
        GeneratePortTreeForIP(ip).then((result) => {
            if (result !== null) {
                portTree = portTree.concat(result);
            }

            loading = false;
        })
            .catch(error => {
                console.error(error);
                loading = false;
            });
    }

</script>
<Row>
    <Column>
        {#each ips as ip}
            <Row>
                <Button kind="secondary" size="small" on:click={() => {
                    ips = [...ips, ""];
                    }}
                        icon={Add}
                        iconDescription="Add">
                </Button>
                <Tile>
                    <TextInput type="text" value="" on:input={(event) => {ip = event.detail.toString()}}/>
                </Tile>
                <Tile>
                    <Button kind="secondary" size="small" on:click={() => {generatePortTree(ip)}}
                            icon={ArrowRight}
                            disabled={ip === ""}
                            iconDescription="Ping">
                    </Button>
                </Tile>
            </Row>
        {/each}

    </Column>
    <Column>
        <Tile>
            {#if loading}
                <StructuredListSkeleton/>
            {:else}
                <TreeView children={portTree} labelText="Open ports"/>
            {/if}
        </Tile>
    </Column>
</Row>