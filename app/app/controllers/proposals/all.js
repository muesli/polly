import Ember from 'ember';

const { inject: { service }, Component } = Ember;

export default Ember.Controller.extend({
  currentUser: service('current-user'),

  activeSmallGrantCount: Ember.computed('smallGrantProposals.@each.moderated', function() {
      var count = 0;
      this.get('smallGrantProposals').forEach(function(entry) {
          if (entry.get('moderated')) {
              count++;
          }
      });
      return count;
  }),
  inactiveSmallGrantCount: Ember.computed('smallGrantProposals.@each.moderated', function() {
      var count = 0;
      this.get('smallGrantProposals').forEach(function(entry) {
          if (!entry.get('moderated')) {
              count++;
          }
      });
      return count;
  }),

  activeLargeGrantCount: Ember.computed('largeGrantProposals.@each.moderated', function() {
      var count = 0;
      this.get('largeGrantProposals').forEach(function(entry) {
          if (entry.get('moderated')) {
              count++;
          }
      });
      return count;
  }),
  inactiveLargeGrantCount: Ember.computed('largeGrantProposals.@each.moderated', function() {
      var count = 0;
      this.get('largeGrantProposals').forEach(function(entry) {
          if (!entry.get('moderated')) {
              count++;
          }
      });
      return count;
  })
});
