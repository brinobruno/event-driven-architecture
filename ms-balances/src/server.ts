import Fastify from 'fastify'
import pg from '@fastify/postgres'
import { kafka } from './external/kafka.js'

const fastify = Fastify({
  logger: true
})

fastify.register(pg, {
  connectionString: 'postgres://postgres:postgres@postgres:5432/wallet'
})

fastify.ready(async () => {
  await fastify.pg.query(`
    CREATE TABLE IF NOT EXISTS balances (
      id SERIAL PRIMARY KEY,
      account_id VARCHAR(255) NOT NULL,
      balance DECIMAL(10, 2) NOT NULL
    )
  `)
  
  try {
    await fastify.pg.query(`
      ALTER TABLE balances ADD CONSTRAINT balances_account_id_unique UNIQUE (account_id)
    `)
  } catch (error: any) {
    // 23505 is the code for unique constraint violation
    // 42P07 is the code for table does not exist
    if (error.code !== '23505' && error.code !== '42P07') {
      throw error
    }
  }
})

fastify.get('/balances/:id', async function handler (request, reply) {
  const { id } = request.params as { id: string }
  const balance = await fastify.pg.query(`
    SELECT * FROM balances WHERE account_id = $1
  `, [id])
  
  if (balance.rows.length === 0) {
    return reply.code(404).send({error: 'Balance not found for account', account_id: id})
  }
  
  return reply.code(200).send({account_id: id, balance: balance.rows[0].balance})
})

const consumer = kafka.consumer({ groupId: 'walletcore' })
await consumer.connect()
await consumer.subscribe({ topic: 'balances', fromBeginning: true })
await consumer.run({
  eachMessage: async ({ topic, partition, message }) => {
    const payload = JSON.parse(message.value?.toString() || '{}')
    const data = payload.Payload
    
    console.info('data', data)

    try {
      await fastify.pg.query(`
        INSERT INTO balances (account_id, balance) 
        VALUES ($1, $2) 
        ON CONFLICT (account_id) 
        DO UPDATE SET balance = EXCLUDED.balance
      `, [data.account_id_from, data.balance_account_id_from])
      
      await fastify.pg.query(`
        INSERT INTO balances (account_id, balance) 
        VALUES ($1, $2) 
        ON CONFLICT (account_id) 
        DO UPDATE SET balance = EXCLUDED.balance
      `, [data.account_id_to, data.balance_account_id_to])
      
      console.log('Updated balances for:', data.account_id_from, 'and', data.account_id_to)
    } catch (error) {
      console.error('Error updating balances:', error)
    }
  }
})

try {
  await fastify.listen({ port: 3003, host: '0.0.0.0' })
} catch (err) {
  fastify.log.error(err)
  process.exit(1)
}