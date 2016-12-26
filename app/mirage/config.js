export default function() {

    this.urlPrefix='http://localhost:9999';
    this.namespace='/v1';

    this.passthrough('/sessions/create');
    this.passthrough('/sessions/signup');
    this.passthrough('/users');
    this.passthrough('/users/:user_id');
    this.passthrough('/proposals');
    this.passthrough('/proposals/:proposal_id');

/*    this.get('/proposals', function() {
      return {
        data:[{
          type: 'proposal',
          id: '1',
          attributes: {
            title: 'Such Proposal',
            recipient: 'Much Recipient',
            enddate: '2016-12-24',
            votes: 20
          }
        },
        {
          type: 'proposal',
          id: '2',
          attributes: {
            title: 'Interesting Title',
            recipient: 'Your Mum',
            enddate: '2017-01-17',
            votes: 1337
          }
        },
        {
          type: 'proposal',
          id: '3',
          attributes: {
            title: 'Fantastic proposal',
            recipient: 'Donkey Kong',
            enddate: '2017-01-02',
            votes: 0
          }
        }]
      };
  });*/

/*    this.post('/invitations');//, 'invitation', 500);
    this.get('/invitations', function() {
      return {
        data:[{
          type: 'invitations',
          id: '1',
          attributes: {
            email: 'muesli@gmail.com'
          }
        },
        {
          type: 'invitations',
          id: '2',
          attributes: {
            email: 'mo@headstrong.de'
          }
        }]
      };
  });*/

  // These comments are here to help you get started. Feel free to delete them.

  /*
    Config (with defaults).

    Note: these only affect routes defined *after* them!
  */

  // this.urlPrefix = '';    // make this `http://localhost:8080`, for example, if your API is on a different server
  // this.namespace = '';    // make this `api`, for example, if your API is namespaced
  // this.timing = 400;      // delay for each request, automatically set to 0 during testing

  /*
    Shorthand cheatsheet:

    this.get('/posts');
    this.post('/posts');
    this.get('/posts/:id');
    this.put('/posts/:id'); // or this.patch
    this.del('/posts/:id');

    http://www.ember-cli-mirage.com/docs/v0.2.x/shorthands/
  */
}
