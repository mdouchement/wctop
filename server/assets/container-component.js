Vue.component('x-container', {
  template: `
  <article class="tile is-vertical is-parent is-2 box notification is-light right-gap">
    <div class="tile is-child">
      <p><strong>
        <span class="icon"><i class="fa fa-laptop"></i></span>
        {{ stats['name'] }}
      </strong></p>
      <p>{{ stats['id'] | trurncateID }}</p>
      <p>{{ stats['started_at'] | runningTime }}</p>
    </div>
    <div class="tile is-child">
      <p>CPU {{ stats['cpu_usage'] }}&#37;
        <progress class="progress is-info" v-bind:value="stats['cpu_usage']" max="100">{{ stats['cpu_usage'] }}&#37;</progress>
      </p>
      <p>RAM {{ stats['mem_usage'] | bytesHumanize }} / {{ stats['mem_limit'] | bytesHumanize }}
        <progress class="progress is-primary" v-bind:value="stats['mem_percent'] +1 " max="100">{{ stats['mem_percent'] }}&#37;</progress>
      </p>
    </div>
  </article>
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
