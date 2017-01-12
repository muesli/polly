import Ember from 'ember';

const { inject: { service } } = Ember;

export default Ember.Controller.extend({
  currentUser: service('current-user'),

  currentUserVoted: Ember.computed('vote.@each.moderated', function() {
      const proposalID = this.get('proposal').get('id');
      var found = false;

      this.get('vote').forEach(function(entry) {
          if (entry.get('proposal').get('id') === proposalID) {
              found = true;
              return;
          }
      });
      return found;
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
    }
  }
});
