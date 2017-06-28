import Ember from 'ember';

const { inject: { service } } = Ember;

export default Ember.Controller.extend({
  currentUser: service('current-user'),

  proposal_has_started: Ember.computed('proposal.starts', function() {
      return this.get('proposal').get('starts') <= new Date();
  }),

  currentUserVoted: Ember.computed('vote.@each.voted', 'proposal', function() {
      const proposalID = this.get('proposal').get('id');
      var found = false;
      var val = 0;

      this.get('vote').forEach(function(entry) {
          if (entry.get('proposal').get('id') === proposalID) {
              found = true;
              if (entry.get('voted')) {
                  val = 1;
              } else {
                  val = -1;
              }
              return;
          }
      });
      return val;
  }),

  actions: {
    moderate(id) {
        this.store.findRecord('proposal', id).then(function(proposal) {
            proposal.set('moderated', true);
            proposal.save();
        });
    },
    vote() {
        const proposalID = this.get('proposal').get('id');
        const newVote = this.store.createRecord('vote', { proposal: this.get('proposal'), voted: true });
        newVote.save().then(
          (/*vote*/) => {
              this.store.findRecord('proposal', proposalID, {reload: true});
          },
          error => {
            alert(error);
          }
        );
    },
    voteDown() {
        const proposalID = this.get('proposal').get('id');
        const newVote = this.store.createRecord('vote', { proposal: this.get('proposal'), voted: false });
        newVote.save().then(
          (/*vote*/) => {
              this.store.findRecord('proposal', proposalID, {reload: true});
          },
          error => {
            alert(error);
          }
        );
    }
  }
});
