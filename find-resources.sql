SELECT config::json->'source'->>'uri' uri, config::json->'source'->>'branch' branch, p.name, r.id
FROM resources r
INNER JOIN pipelines p ON p.id = r.pipeline_id
WHERE r.active=TRUE 
    AND r.name='git'
    AND r.paused=FALSE 
    AND p.paused=FALSE 
    AND r.pipeline_id = p.id 
    AND config::json->'source'->>'uri' IN  ('git@github.com:springernature/oscar.git', 'git@github.com:springernature/oscar')
    AND config::json->'source'->>'branch' = 'master'

