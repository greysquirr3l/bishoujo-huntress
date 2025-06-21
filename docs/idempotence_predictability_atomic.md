# Idempotence, Predictability, and Atomic Actions in Go

Idempotence is a critical property in robust software systems, particularly when working with distributed architectures, APIs, and systems that may experience failures and retries. Our projects follow DDD, CQRS, and Clean Architecture principles, which all benefit from idempotent operations.

## What is Idempotence?

- An operation is idempotent when applying it multiple times produces the same result as applying it once
- Repeated executions don't cause additional side effects beyond the first execution
- Examples include HTTP methods like GET, PUT, and DELETE (when properly implemented)
- In our domain models, idempotent methods ensure consistent state transitions

### Why Idempotence Matters

- **Resilience to Failures**: Operations can be safely retried without causing duplicate effects
- **Distributed Systems**: Makes systems more reliable when messages might be delivered more than once
- **API Consistency**: Provides predictable behavior for consumers of your services
- **Improved Testability**: Idempotent operations are easier to test and reason about
- **Better DDD Implementation**: Ensures domain operations produce consistent, reliable results

### Implementing Idempotence in the Context of Clean Architecture and CQRS

The principles in our project guidelines document highlight the importance of separating concerns through DDD, CQRS, and Clean Architecture. Idempotence strengthens these patterns:

1. **Use Idempotency Keys with Command Handlers**:

   This aligns with our CQRS pattern from the project guidelines:

   ```go
   func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (string, error) {
       // Check if command was already processed using idempotency key
       userID, processed, err := h.idempotencyStore.CheckProcessed(ctx, cmd.IdempotencyKey)
       if err != nil {
           return "", err
       }
       if processed {
           return userID, nil // Return previous result
       }

       // Domain validation within command handler (as shown in project guidelines)
       email, err := domain.NewEmail(cmd.Email)
       if err != nil {
           return "", err
       }
       password, err := domain.NewPassword(cmd.Password)
       if err != nil {
           return "", err
       }

       // Create domain entity
       user := domain.NewUser(email, password)

       // Persist through repository
       if err := h.userRepo.Save(ctx, user); err != nil {
           return "", err
       }

       // Store result with idempotency key
       if err := h.idempotencyStore.StoreResult(ctx, cmd.IdempotencyKey, user.ID.String()); err != nil {
           return "", err
       }

       return user.ID.String(), nil
   }
   ```

2. **Use Conditional Operations in Repositories**:

   This implements the Repository pattern described in our project guidelines:

   ```go
   // Idempotent repository save implementation in the infrastructure layer
   func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
       // Use upsert pattern for idempotence
       _, err := r.db.ExecContext(ctx,
           `INSERT INTO users (id, email, password_hash, role, active)
            VALUES ($1, $2, $3, $4, $5)
            ON CONFLICT (id) DO UPDATE
            SET email = $2, password_hash = $3, role = $4, active = $5`,
           user.ID.String(), string(user.Email), string(user.Password),
           string(user.Role), user.Active)
       return err
   }
   ```

3. **Domain Event Idempotence**:

   Our project guidelines emphasize domain events; here's how to make them idempotent:

   ```go
   // In domain layer: Event with deduplication info
   type UserCreatedEvent struct {
       EventID   string    // Unique event ID for deduplication
       UserID    string
       Email     string
       OccurredAt time.Time
   }

   // In infrastructure layer: Event publisher with deduplication
   func (p *KafkaEventPublisher) PublishUserCreated(ctx context.Context, event domain.UserCreatedEvent) error {
       // Check if this event was already published
       if published, _ := p.eventStore.WasPublished(event.EventID); published {
           return nil // Event already published, skip
       }

       // Publish event...

       // Mark as published
       return p.eventStore.MarkPublished(event.EventID)
   }
   ```

## Ensuring Predictability

Predictability complements idempotence by ensuring that:

- **Deterministic Behavior**: Given the same inputs, operations always produce the same outputs
- **Explicit Side Effects**: All side effects are documented, understood, and consistent
- **Proper Error Handling**: Errors are handled predictably and don't leave the system in an inconsistent state
- **Avoid Hidden State**: State changes should be explicit and traceable

### Example: Deterministic Domain Methods

From our project guidelines, we can enhance our domain methods to be predictable:

```go
// Instead of this:
func (t *Training) CanBeAttendedBy(u *User) error {
    // Unpredictable: might depend on current time
    if time.Now().After(t.EndTime) {
        return ErrTrainingEnded
    }
    // Rest of logic...
}

// Do this:
func (t *Training) CanBeAttendedByAt(u *User, checkTime time.Time) error {
    // Predictable: explicitly depends on provided time
    if checkTime.After(t.EndTime) {
        return ErrTrainingEnded
    }
    // Rest of logic...
}
```

## Combining with DDD & CQRS

When combined with Domain-Driven Design and CQRS:

- **Commands should be idempotent when possible**: Command handlers should detect and handle duplicate commands
- **Queries are naturally idempotent**: They don't change state, making them safe to retry
- **Domain events should capture state transitions in a predictable manner**: Include all necessary context for consumers
- **Repositories should implement idempotent save operations**: Use conditional updates or optimistic concurrency control

## Atomic actions are preferred

In our project guidelines, we emphasize clean architecture and domain-driven design. Atomic operations play a crucial role in maintaining the integrity of our domain model:

- An atomic action in programming means that a sequence of operations is treated as a single, indivisible unit
- No intermediate states are visible to other processes or threads during the execution
- Either the entire action is completed, or none of it is, ensuring data consistency

### Implementing Atomic Operations in Clean Architecture

Within our layer architecture, we should implement atomic operations at the appropriate levels:

1. **Domain Layer**: Model invariants that must be maintained atomically

   ```go
   // From project guidelines: Domain logic within entity
   func (u *User) ChangePassword(currentPassword, newPassword Password) error {
       // This entire operation is conceptually atomic
       if !u.Password.Matches(currentPassword) {
           return ErrInvalidPassword
       }
       u.Password = newPassword
       return nil
   }
   ```

2. **Application Layer**: Use transactions for operations that span multiple aggregates

   ```go
   // Command Handler with transaction
   func (h *TransferSubscriptionHandler) Handle(ctx context.Context, cmd TransferSubscriptionCommand) error {
       // Begin transaction to ensure atomicity
       tx, err := h.db.BeginTx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()

       userRepo := h.userRepo.WithTx(tx)

       // Load both users
       fromUser, err := userRepo.GetByID(ctx, cmd.FromUserID)
       if err != nil {
           return err
       }

       toUser, err := userRepo.GetByID(ctx, cmd.ToUserID)
       if err != nil {
           return err
       }

       // Domain operations
       subscription, err := fromUser.RemoveSubscription()
       if err != nil {
           return err
       }

       if err := toUser.AddSubscription(subscription); err != nil {
           return err
       }

       // Save both users atomically
       if err := userRepo.Save(ctx, fromUser); err != nil {
           return err
       }
       if err := userRepo.Save(ctx, toUser); err != nil {
           return err
       }

       // Commit the transaction to ensure atomicity
       return tx.Commit()
   }
   ```

3. **Infrastructure Layer**: Leverage database transactions to ensure atomicity

   ```go
   // Infrastructure implementation of WithTx from project guidelines
   func (r *SQLUserRepository) WithTx(tx *sql.Tx) UserRepository {
       // Creates a new repository instance that uses the transaction
       return &SQLUserRepository{db: tx}
   }
   ```

### Atomicity and CQRS

The Command Query Responsibility Segregation pattern in our project guidelines naturally supports atomic operations:

- **Commands**: Should be atomic operations that either fully succeed or fail
- **Queries**: Don't modify state, so atomicity concerns don't apply
- **Event Publishing**: Should be part of the same atomic operation as state changes

### Best Practices for Atomic Operations

1. **Identify Transactional Boundaries**: Determine which operations must succeed or fail together
2. **Use Database Transactions**: For persistence operations that span multiple entities/tables
3. **Keep Transactions Short**: Long-running transactions can cause contention
4. **Handle Concurrent Modifications**: Use optimistic or pessimistic concurrency control
5. **Consider Eventual Consistency**: For operations that span service boundaries

### Testing Atomic Operations

When testing, verify that operations are truly atomic:

```go
// Test that transfer is atomic - either both users are updated or neither is
func TestTransferSubscriptionAtomicity(t *testing.T) {
    // Setup test
    handler := NewTransferSubscriptionHandler(mockRepo, mockDB)

    // Make the save of second user fail
    mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).
        Return(nil).Once().
        Return(errors.New("database error")).Once()

    // Execute command
    err := handler.Handle(context.Background(), TransferSubscriptionCommand{
        FromUserID: "user1",
        ToUserID: "user2",
    })

    // Assert transaction rolled back
    assert.Error(t, err)

    // Verify neither user was changed when we fetch them again
    fromUser, _ := realRepo.GetByID(context.Background(), "user1")
    toUser, _ := realRepo.GetByID(context.Background(), "user2")

    assert.True(t, fromUser.HasSubscription())
    assert.False(t, toUser.HasSubscription())
}
```
