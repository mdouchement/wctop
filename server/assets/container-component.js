Vue.component('x-container', {
  template: `
  <div class="column x-is-container">
    <article class="tile is-vertical box notification is-light no-padding">
      <div class="tile is-parent">
        <div class="tile is-child">
          <p><strong>
            <span class="icon"><i class="fa fa-laptop"></i></span>
            {{ stats['name'] }}
          </strong></p>
          <p>{{ stats['id'] | trurncateID }}</p>
          <p>{{ stats['started_at'] | runningTime }}</p>
        </div>
      </div>
      <div class="tile is-parent x-small-font">
        <div class="tile is-child">
          <p>CPU {{ stats['cpu_usage'] }}&#37;
            <vue-simple-progress bar-color="#209CEE" bg-color="#DBDBDB" size="big" v-bind:val="stats['cpu_usage']"></vue-simple-progress>
          </p>
          <p>RAM {{ stats['mem_usage'] | bytesHumanize }} / {{ stats['mem_limit'] | bytesHumanize }}
            <vue-simple-progress bar-color="#00D1B2" bg-color="#DBDBDB" size="big" v-bind:val="stats['mem_percent'] +1"></vue-simple-progress>
          </p>
        </div>
      </div>
      <div class="tile is-parent x-small-font">
        <div class="tile is-child is-vertical">
          <p>Tx: {{ stats['net_tx'] | bytesHumanize }}</p>
          <p>Rx: {{ stats['net_rx'] | bytesHumanize }}</p>
        </div>
        <div class="tile is-child">
          <p>Read: {{ stats['io_bytes_read'] | bytesHumanize }}</p>
          <p>Write: {{ stats['io_bytes_write'] | bytesHumanize }}</p>
        </div>
      </div>
    </article>
  </div>
  `,
  props: ['stats'],
  filters: {
    trurncateID: function (value) {
      return value.substring(0, 16) + '...'
    },
    bytesHumanize: function (value) {
      return bytesHumanize(value);
    },
    runningTime: function (date) {
      return moment(date).fromNow();
    }
  }
})
