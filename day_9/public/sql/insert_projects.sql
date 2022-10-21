INSERT INTO public.tb_projects(
	name, start_date, end_date, description, technologies, image)
	VALUES ('percobaan makan', NOW()::date, '2023-05-15', 'melakukan hal hal aneh', ARRAY['alpha', 'beta'],'?'),('momen aneh terjadi', CURRENT_DATE, '2024-10-05', 'melakukan hal hal percobaan',ARRAY['beta','test','gen'],'?');